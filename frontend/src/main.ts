import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router'
import Toast, { type PluginOptions } from "vue-toastification"
import "vue-toastification/dist/index.css"

const app = createApp(App)

const toastOptions: PluginOptions = {
  position: "top-right",
  timeout: 4000,
  closeOnClick: true,
  pauseOnFocusLoss: true,
  pauseOnHover: true,
  draggable: true,
  draggablePercent: 0.6,
  showCloseButtonOnHover: false,
  hideProgressBar: false,
  closeButton: "button",
  icon: true,
  rtl: false
}
app.use(Toast, toastOptions)

app.use(createPinia())
app.use(router)

app.mount('#app')
