package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8smanager-demo/config"
	"log"
	"net/http"
	"time"
)

var Terminal terminal

type terminal struct{}

// 消息内容
type terminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

// TerminalSession 交互的结构体，接管输入和输出
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

// 初始化一个websocket.Upgrader类的对象，用于http协议升级为
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader
}()

// WsHandler 定义websocket的handler方法
func (t *terminal) WsHandler(w http.ResponseWriter, r *http.Request) {
	//加载k8s配置
	conf, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		fmt.Println("创建k8s配置失败, " + err.Error())
	}
	//解析form入参，获取namespace、podName、containerName参数
	if err := r.ParseForm(); err != nil {
		return
	}
	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	fmt.Println("exec pod: %s, container: %s, namespace: %s\n", podName,
		containerName, namespace)
	//new一个TerminalSession类型的pty实例
	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		fmt.Println("get pty failed: %v\n", err)
		return
	}
	//处理关闭
	defer func() {
		fmt.Println("close session.")
		pty.Close()
	}()
	// 初始化pod所在的corev1资源组
	// PodExecOptions struct 包括Container stdout stdout Command 等结构
	// scheme.ParameterCodec 应该是pod 的GVK （GroupVersion & Kind）之类的
	// URL长相:
	// https://192.168.1.11:6443/api/v1/namespaces/default/pods/nginx-wf2-778d88d7c7rmsk/exec?command=%2Fbin%2Fbash&container=nginxwf2&stderr=true&stdin=true&stdout=true&tty=true
	req := K8s.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command:   []string{"/bin/bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	fmt.Println(req.URL())
	//remotecommand 主要实现了http 转 SPDY 添加X-Stream-Protocol-Version相关header 并发送请 求
	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		return
	}
	// 建立链接之后从请求的sream中发送、读取数据
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	})
	if err != nil {
		msg := fmt.Sprintf("Exec to pod error! err: %v", err)
		fmt.Println(msg)
		//将报错返回出去
		pty.Write([]byte(msg))
		//标记退出stream流
		pty.Done()
	}
}

// NewTerminalSession 创建TerminalSession类型的对象并返回
func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	// 升级ws协议
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, errors.New("升级websocket失败" + err.Error())
	}
	// new
	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	// 返回
	return session, nil
}

// Done 标记关闭的方法
func (t *TerminalSession) Done() {
	close(t.doneChan)
}

// Close 关闭的方法
func (t *TerminalSession) Close() error {
	return t.wsConn.Close()
}

// Next resize方法，以及是否退出终端
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

// 读数据的方法
// 返回值int是都读成功了多少数据
func (t *TerminalSession) Read(p []byte) (int, error) {
	//从ws中读取消息
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		fmt.Printf("read message err: %v", err)
		return 0, err
	}
	// 反序列化
	var msg terminalMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		fmt.Printf("read message err: %v", err)
		return 0, err
	}
	// 逻辑判断
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{
			Width:  msg.Cols,
			Height: msg.Rows,
		}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		fmt.Printf("unknown message type '%s'", msg.Operation)
		return 0, errors.New("unknown message type '%s'" + msg.Operation)
	}
}

// 写数据的方法
func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(terminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		log.Printf("write parse message err: %v", err)
		return 0, err
	}
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("write message err: %v", err)
		return 0, err
	}
	return len(p), nil
}
