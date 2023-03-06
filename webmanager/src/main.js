import { createApp } from 'vue'
import App from './App.vue'
//引入element plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import router from './router'

// 引入图标
import * as ELIcons from '@element-plus/icons-vue'


const app = createApp(App)

app.use(router)


app.use(ElementPlus)


for (let iconName in ELIcons) {
    app.component(iconName, ELIcons[iconName])
}


app.mount('#app')
