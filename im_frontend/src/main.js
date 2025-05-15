import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import axios from 'axios';

axios.defaults.withCredentials = true;
axios.defaults.baseURL = 'http://localhost:8081';
//axios.defaults.baseURL = 'http://192.168.71.154:8080';

//添加响应拦截器
axios.interceptors.response.use(
  (response) => {
    // 对响应数据做点什么
    return response;
  },
  (error) => {
    // 对响应错误做点什么
    if (error.response && (error.response.status === 401 || error.response.status === 403)) {
      // 如果后端返回 401 或 403 状态码，表示未认证或认证过期
      router.push({ name: 'Login' });
    }
    return Promise.reject(error);
  }
);

createApp(App).use(router).use(store).mount('#app');