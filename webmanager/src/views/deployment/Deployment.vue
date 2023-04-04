<template>
  <div class="deploy">
    <el-row>
      <!-- 头部1 -->
      <el-col :span="24">
        <div>
          <el-card class="deploy-head-card" shadow="never" :body-style="{padding:'10px'}">
            <el-row>
<!--              命名空间-->
              <el-col :span="6">
              <div>
                <span>命名空间:</span>
                <el-select v-model="namespaceValue" filterable placeholder="Select">
                  <el-option
                      v-for="(item,index) in namespaceList"
                      :key="index"
                      :label="item.metadata.name"
                      :value="item.metadata.name"
                  >
                  </el-option>
                </el-select>
              </div>
              </el-col>
              <!--              刷新按钮-->
              <el-col :span="2" :offset="16">
                <dev>
                  <el-button style="border-radius:2px" icon="Refresh" plain>刷新</el-button>
                </dev>
              </el-col>
              <el-col>

              </el-col>
            </el-row>
          </el-card>
        </div>
      </el-col>
      <!-- 头部2 -->
      <el-col :span="24">
        <div>
          <el-card class="deploy-head-card" shadow="never" :body-style="{padding:'10px'}">
            <el-row>
<!--              创建按钮-->
              <el-col :span="6">
                <div>
                  <el-button
                      style="border-radius:2px;"
                      icon="Edit"
                      type="primary"
                      v-loading.fullscreen.lock="fullscreenLoading"
                      @click="createDeploymentDrawer=true">
                    创建
                  </el-button>
                </div>
              </el-col>
<!--              搜索框和搜索按钮-->
              <el-col :span="6">
                <div>
                  <el-input class="deploy-head-search" clearable placeholder="请输入" v-model="searchInput"></el-input>
                  <el-button style="border-radius:2px;"
                    icon="Search" type="primary">
                    搜索
                  </el-button>
                </div>
              </el-col>
            </el-row>
          </el-card>

        </div>
      </el-col>
      <!-- 数据表格 -->
      <el-col :span="24"></el-col>
    </el-row>
    <!--  https://element-plus.org/zh-CN/component/drawer.html  -->
    <el-drawer
    v-model="createDeploymentDrawer"
    :direction="direction"
    :before-close="handleClose" >
<!--      标题-->
      <template #title>
        <h4>创建Deployment</h4>
      </template>
<!--      body,填写表单属性-->
      <template #default>
        <el-row type="flex" justify="center" >
          <el-col :span="20">
            <el-form
                ref="createDeployment"
                :rules="createDeploymentRules"
                :model="createDeployment"
                label-width="80px"
            >
              <el-form-item class="deploy-create-form" label="名称" prop="name">
                <el-input v-model="createDeployment.name"></el-input>
              </el-form-item>
              <el-form-item class="deploy-create-form" label="命名空间" prop="namespace">
                <el-select v-model="createDeployment.namespace" filterable placeholder="Select">
                  <el-option
                      v-for="(item,index) in namespaceList"
                      :key="index"
                      :label="item.metadata.name"
                      :value="item.metadata.name"
                  >
                  </el-option>
                </el-select>
              </el-form-item>
              <el-form-item class="deploy-create-form" label="副本数" prop="replicas">
                <el-input-number v-model="createDeployment.replicas" :min="1" :max="10">
                </el-input-number>
                <el-popover
                    placement="top"
                    :width="100"
                    trigger="hover"
                    content="申请副本数上限为10个"
                >
                  <template #reference>
                    <el-icon style="width:2em;font-size: 18px; color: #4796EE"><WarningFilled /></el-icon>
                  </template>
                </el-popover>
              </el-form-item>
              <el-form-item class="deploy-create-form" label="镜像" prop="image">
                <el-input v-model="createDeployment.image"></el-input>
              </el-form-item>
              <el-form-item class="deploy-create-form" label="标签" prop="label_str">
                <el-input v-model="createDeployment.label_str" placeholder="示例: project=ms,app=gateway"></el-input>
              </el-form-item>
              <el-form-item class="deploy-create-form" label="资源额配" prop="resource">
                <el-select v-model="createDeployment.resource" placeholder="请选择">
                  <el-option value="0.5/1" label="0.5C1G"></el-option>
                    <el-option value="1/2" label="1C2G"></el-option>
                    <el-option value="2/4" label="2C4G"></el-option>
                    <el-option value="4/8" label="4C8G"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item class="deploy-create-form" label="容器端口" prop="container_port">
                <el-input v-model="createDeployment.container_port" placeholder="示例: 80"></el-input>
              </el-form-item>
              <el-form-item class="deploy-create-form" label="健康检查">
                <el-switch v-model="createDeployment.health_check" />
              </el-form-item>
              <el-form-item class="deploy-create-form" label="检查路径">
                <el-input v-model="createDeployment.health_path" placeholder="示例: /health"></el-input>
              </el-form-item>
            </el-form>
          </el-col>
        </el-row>
      </template>
<!--      footer,处理提交和取消-->
      <template #footer>
        <el-button @click="createDeploymentDrawer = false">取消</el-button>
        <el-button type="primary" @click="submiForm('createDeployment')">立即创建</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script>
import common from "../../common/Config"
import httpClient from "@/utils/request";
export default {
  name: "Deployment",
  data() {
    return {
      // 搜索
      searchInput: "",
      // 命名空间
      namespaceValue: 'default',
      namespaceList: [],
      namespaceListUrl: common.k8sNamespaceList,
      // 创建
      createDeploymentDrawer: false,
      fullscreenLoading: false,
      direction: 'rtl',
      // 创建deployment的属性
      createDeployment: {
        name:'',
        namespace:'',
        replicas:'',
        image: '',
        resource:'',
        health_check: false,
        health_path: '',
        label_str: '',
        label:{},
      },
      // 发送请求时的参数
      createDeploymentData: {
        url: common.k8sDeploymentCreate,
        params: {

        },
      },
      // 创建deployment表单的校验规则
      createDeploymentRules: {
        name: [{
          required: true,
          message: '请填写名称',
          trigger: 'change'
        }],
        image: [{
          required: true,
          message: '请填写镜像',
          trigger: 'change'
        }],
        namespace: [{
          required: true,
          message: '请选择命名空间',
          trigger: 'change'
        }],
        resource: [{
          required: true,
          message: '请选择配额',
          trigger: 'change'
        }],
        label_str: [{
          required: true,
          message: '请填写标签',
          trigger: 'change'
        }],
        container_port: [{
          required: true,
          message: '请填写容器端口',
          trigger: 'change'
        }],
      },
    }
  },
  methods: {

    // 处理抽屉的关闭，double check增加体验
    handleClose(done) {
      this.$confirm('确认关闭')
          .then(() => {
            done();
          })
          .catch(() => {

          })
    },
    getNamespaces() {
      httpClient.get(this.namespaceListUrl).then(res => {
        this.namespaceList = res.data
      }).catch(res => {
        this.$message.error({
          message: res.msg,
        })
      })
    },
    // 创建deployment
    submiForm(formName) {
      this.$refs[formName].validate((valid)=> {
        if (valid) {
          this.createDeployFunc()
        } else {
          return false
        }
      })
    },
    createDeployFunc() {
      // 正则匹配验证label
      let reg = new Regexp("(^[A-Za-z]+=[A-Za-z0-9]+).*")
      // 负载均衡资源54.13

      // 处理label,cpu和内存

      // 赋值

      // 发起请求
    },
  },
  watch: {
    // 监听namespace的值,若发生变化，则执行handler方法中的内容
    namespaceValue : {
      handler() {
        // 将namespace的值存入本地，用于path切换时依旧能获取到
        localStorage.setItem('namespace',this.namespaceValue)
        // console.log(this.namespaceValue)
      }
    }
  },
  beforeMount() {
    // 加载页面的时候先获取localStorage中的namespace值，若获取不到则默认default
    if (localStorage.getItem('namespace')!==undefined && localStorage.getItem('namespace')!=null){
      this.namespaceValue=localStorage.getItem('namespace')
    }
    this.getNamespaces()
  },
}
</script>

<style scoped>
.deploy-head-card {
  border-right: 1px;
  margin: 5px;
}
/*搜索框*/
.deploy-head-search{
  width: 160px;
}
</style>