import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import baseStore from './components/base/store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import '@/assets/icon-font/iconfont.css'
import '@/assets/icon-font/iconfont.css'
import 'viewerjs/dist/viewer.css'
import VueViewer from 'v-viewer'
import i18n from './i18n.js'
// import VConsole from 'vconsole'
// const vConsole = new VConsole()

const app = createApp(App)

// register base layout modules under their own namespace if missing
if (!store.hasModule('layout')) {
  store.registerModule('layout', baseStore.state ? baseStore.state.layout || baseStore._modulesNamespaceMap['layout/']?.context?.module : baseStore)
}
// if (!store.hasModule('dashboard')) {
//   store.registerModule('dashboard', baseStore.state ? baseStore.state.dashboard || baseStore._modulesNamespaceMap['dashboard/']?.context?.module : baseStore)
// }

app.config.productionTip = false;
window.log = console.log;

app.use(ElementPlus);
app.use(VueViewer);
app.use(store);
app.use(router);
app.use(i18n);


app.mount('#app');