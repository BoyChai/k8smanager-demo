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
      createDeploymentDrawer: false
    }
  },
  methods: {
    getNamespaces() {
      httpClient.get(this.namespaceListUrl).then(res => {
        this.namespaceList = res.data
      }).catch(res => {
        this.$message.error({
          message: res.msg,
        })
      })
    }
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