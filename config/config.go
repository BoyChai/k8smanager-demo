package config

import "time"

const (
	ListenAddr = "0.0.0.0:9090"
	KubeConfig = "C:\\Users\\BoyChai\\.kube\\config"
	// PodLogTailLine tail的日志行数
	//tail -n 2000
	PodLogTailLine = 2000

	// 数据库配置

	DbHost = "host"
	DbPort = 3306
	DbName = "name"
	DbUser = "user"
	DbPass = "pass"

	// 连接池的配置

	MaxIdleConns = 10               // 最大空闲链接
	MaxOpenConns = 100              // 最大连接数
	MaxLifeTime  = 30 * time.Second // 最大存活时间

)
