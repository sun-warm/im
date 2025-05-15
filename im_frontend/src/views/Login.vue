// filepath: /Users/sunwen/Projects/FrontEnd/frontend/src/views/Login.vue
<template>
  <div class="login-container">
    <h2>登录</h2>
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="username">用户名</label>
        <input type="text" id="username" v-model="username" required />
      </div>
      <div class="form-group">
        <label for="password">密码</label>
        <input type="password" id="password" v-model="password" required />
      </div>
      <button type="submit">登录</button>
    </form>
  </div>
</template>

<script>
import axios from 'axios';
export default {
  name: 'LoginPage',
  data() {
    return {
      username: '',
      password: ''
    };
  },
  mounted() {
    console.log('LoginPage component has been mounted');
  },
  methods: {
    async handleLogin() {
      // 捕获用户输入的账号密码
      console.log('用户名:', this.username);
      console.log('密码:', this.password);
      // 发送登陆请求
      try {
        const response = await axios(
          {
            method: 'post',
            url: '/login',
            headers: {
              'Content-Type': 'application/json',
            },
            data: {
              user_name: this.username,
              pass_word: this.password
            },
            withCredentials: true
          });
        console.log('登录成功:', response.data);
        
        sessionStorage.setItem('user_name', this.username);
        console.log('user_name:', this.username);
        
        // 模拟登录成功，跳转到 Home 页面
        this.$router.push({ name: 'Chat' });
      } catch (error) {
        console.error('登录失败:', error);
      }
    }
  }
};
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
}

.form-group input {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

button {
  width: 100%;
  padding: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #38a169;
}
</style>
