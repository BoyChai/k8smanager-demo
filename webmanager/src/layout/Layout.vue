<template>
  <div class="common-layout">
<!--    container布局-->
    <el-container style="height: 100vh">
<!--      侧边栏导航栏-->
      <el-aside class="aside" :width=asideWidth>
        <el-affix class="aside-logo">
          <el-image class="logo-image" :src="logo" />
          <span :class="[isCollapse?'is-collapse':'']">
            <span class="logo-name" v-if="!isCollapse">Kubernetes</span>
          </span>
        </el-affix>
<!--        定义vue router模式，跟路由规则中的path绑定-->
<!--        default-active 默认激活的菜单栏，这里根据打开的path来找到对应的栏-->
        <el-menu class="aside-menu"
                  router
                 :default-active="$route.path"
                 :collapse="isCollapse"
                 background-color="#131b27"
                 text-color="#bfcbd9"
                 active-text-color="20a0ff">
<!--          routers就是router/index.js中的routes-->
          <div v-for="menu in routers" :key="menu">
<!--            第一种情况，路由规则children只有1个的菜单栏-->
            <el-menu-item class="aside-menu-item" v-if="menu.children && menu.children.length == 1" :index="menu.children[0].path">
<!--              处理图标和菜单栏的名字-->
              <el-icon><component :is="menu.children[0].icon"></component></el-icon>
              <template #title>
                {{menu.children[0].name}}
              </template>
            </el-menu-item>
<!--            第二种情况，路由规则children大于1个的菜单栏-->
            <el-sub-menu class="aside-submenu" v-else-if="menu.children && menu.children.length > 1" :index="menu.path">
<!--              处理父菜单栏-->
              <template #title>
                <el-icon><component :is="menu.icon"></component></el-icon>
                <span :class="[isCollapse?'is-collapse':'']">{{menu.name}}</span>
<!--                    <span v-if="!isCollapse">{{menu.name}}</span>-->
                  </template>
<!--              处理子菜单栏-->
              <el-menu-item class="aside-childitem" v-for="child in menu.children" :key="child" :index="child.path">
                <template #title>
                  {{child.name}}
                </template>
              </el-menu-item>
            </el-sub-menu>
          </div>
        </el-menu>
      </el-aside>
<!--      -->
      <el-container>
        <el-header class="header">
          <el-row :gutter="20">
<!--            折叠按钮-->
            <el-col :span="1">
              <div class="header-collapse">
<!--                通过isCollapse来控制状态-->
                <el-icon @click="onCollapse"><component :is="isCollapse?'expand':'fold'"></component></el-icon>
              </div>
            </el-col>
<!--            面包屑-->
            <el-col :span="10">
              <div class="header-breadcrumb">
                <el-breadcrumb separator="/">
<!--                 最外层的工作台 写死的-->
                  <el-breadcrumb-item :to="{path:'/'}">工作台</el-breadcrumb-item>
<!--                  循环出路径的面包屑-->
                  <template v-for="(matched,m) in this.$route.matched" :key="m">
                    <el-breadcrumb-item v-if="matched.name != undefined">
                      {{matched.name}}
                    </el-breadcrumb-item>
                  </template>
                </el-breadcrumb>
              </div>
            </el-col>
<!--            用户信息-->
            <el-col></el-col>
          </el-row>
        </el-header>
        <el-main>main</el-main>
        <el-footer>footer</el-footer>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import {useRouter} from "vue-router";

export default {
  name: "Layout",
  data() {
    return {
      logo:require('@/assets/img/k8s-metrics.png'),
      asideWidth:'220px',
      isCollapse:false,
      routers: []
    }
  },
  beforeMount() {
    //  拿到router对象
    this.routers = useRouter().options.routes
  },
  methods: {
    onCollapse() {
      console.log(this.asideWidth)
      // true为折叠状态
      if (!this.isCollapse){
        // 缩减
        this.asideWidth='64px'
      } else {
        // 展开
        this.asideWidth='220px'
      }
      this.isCollapse = !this.isCollapse
    }
  }

}
</script>

<style scoped>
/*aside属性*/
  .aside {
    transition: all .5s;
    background-color: #131b27;
  }
  .aside-logo{
    z-index: 1200;
    color: aliceblue;
    height: 60px;
  }
  .logo-image{
    width: 40px;
    height: 40px;
    top: 12px;
    padding-left: 12px;
  }
  .logo-name{
    font-size: 20px;
    font-weight: bold;
    padding:10px;
  }
  .is-collapse{
    display: none;
  }
  /*清除滚动轴*/
  .aside::-webkit-scrollbar{
    display: none;
  }
  /*清除边框左右宽度*/
  .aside-menu{
    border-right-width: 0;
  }

/*  菜单栏背景色*/
  .aside-menu-item.is-active{
    background-color: #1f2a3a;
  }
  .aside-menu-item:hover{
    background-color: #142c4e;
  }

  .aside-childitem.is-active {
    background-color: #1f2a3a;
  }
  .aside-childitem:hover{
    background-color: #142c4e;
  }
/*  header的css属性*/
.header{
  z-index: 1200;
  line-height: 60px;
  font-size: 24px;
  box-shadow: 0 2px 4px rgba(0,0,0 ,.12),0 0 6px rgba(0,0,0,.04);
}
.header-collapse{
  cursor: pointer;
}
/* 面包屑 */
.header-breadcrumb{
  padding-top: 0.9em;
}
</style>