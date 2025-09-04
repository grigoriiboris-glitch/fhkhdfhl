import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import '@/assets/icon-font/iconfont.css'
import '@/assets/icon-font/iconfont.css'
import 'viewerjs/dist/viewer.css'
import VueViewer from 'v-viewer'
import i18n from './i18n.js'

const app = createApp(App)

app.config.productionTip = false;
window.log = console.log;

app.use(ElementPlus);
app.use(VueViewer);
app.use(store);
app.use(router);
app.use(i18n);

app.mount('#app');

import icon from '@/components/base/components/ui/icon';