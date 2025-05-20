// filepath: /Users/sunwen/Projects/FrontEnd/frontend/src/views/Login.vue
<template>
  <div class="auth-container">
    <!-- 登录表单 -->
    <div v-if="!showRegister" class="auth-card login-card">
      <h2>登录</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">用户名</label>
          <input
            type="text"
            id="username"
            v-model="username"
            placeholder="请输入用户名"
            required
          />
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input
            type="password"
            id="password"
            v-model="password"
            placeholder="请输入密码"
            required
          />
        </div>
        <button type="submit" class="btn">登录</button>
      </form>
      <div class="switch-link">
        <span>还没有账号？</span>
        <button @click="showRegister = true" class="btn-link">注册</button>
      </div>
    </div>

    <!-- 注册表单 -->
    <div v-if="showRegister" class="auth-card register-card">
      <h2>注册</h2>
      <form @submit.prevent="register">
        <div class="form-group">
          <label for="register-UserID">昵称</label>
          <input
            type="userid"
            id="register-userid"
            v-model="registerUserID"
            placeholder="请输入昵称"
            required
          />
        </div>
        <div class="form-group">
          <label for="register-username">用户名</label>
          <input
            type="text"
            id="register-username"
            v-model="registerUsername"
            placeholder="请输入用户名"
            required
          />
        </div>
        <div class="form-group">
          <label for="register-password">密码</label>
          <input
            type="password"
            id="register-password"
            v-model="registerPassword"
            placeholder="请输入密码"
            required
          />
        </div>
        <div class="form-group">
          <label for="confirm-password">确认密码</label>
          <input
            type="password"
            id="confirm-password"
            v-model="confirmPassword"
            placeholder="请再次输入密码"
            required
          />
        </div>
        <button type="submit" class="btn">注册</button>
        <div class="switch-link">
          <span>已有账号？</span>
          <button @click="showRegister = false" class="btn-link">登录</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "LoginPage",
  data() {
    return {
      username: "",
      password: "",
      showRegister: false, // 控制注册表单显示
      registerUserID:"",
      registerUsername: "",
      registerPassword: "",
      confirmPassword: "",
    };
  },
  methods: {
    async handleLogin() {
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
              pass_word: this.password,
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
    },
    async register() {
      if (this.registerPassword !== this.confirmPassword) {
        alert("两次输入的密码不一致");
        return;
      }
      try {
        const response = await axios(
          {
            method: 'post',
            url: '/register',
            headers: {
              'Content-Type': 'application/json',
            },
            data: {
              user_name: this.registerUsername,
              pass_word: this.registerPassword,
              user_id: this.registerUserID,
            },
          });
        console.log('注册成功:', response.data);
        alert("注册成功");
        this.showRegister = false // 注册成功后切换回登录界面
      } catch (error) {
        console.error('注册失败:' +  error);
        alert("注册失败", error);
      }
    },
  },
};
</script>


<style scoped>
html,
body {
  margin: 0;
  padding: 0;
  height: 100%;
  box-sizing: border-box;
}

.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh; /* 占满整个视口高度 */
  background-color: #f5f5f5;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 20px;
  border-radius: 8px;
  background-color: #fff;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, opacity 0.3s ease;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
  font-size: 24px;
  color: #333;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #555;
}

.form-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-sizing: border-box;
}

.btn {
  width: 100%;
  padding: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.btn:hover {
  background-color: #369f6e;
}

.btn-link {
  background: none;
  border: none;
  color: #42b983;
  cursor: pointer;
  text-decoration: underline;
  font-size: 14px;
}

.btn-link:hover {
  color: #369f6e;
}

.switch-link {
  text-align: center;
  margin-top: 10px;
  font-size: 14px;
  color: #555;
}
</style>