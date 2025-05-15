
<template>
  <div class="chat-container">
    <!-- 左侧聊天对象列表 -->
    <div class="chat-list">
      <div class="search-bar">
        <input
          type="text"
          v-model="searchQuery"
          placeholder="搜索用户名..."
          @input="searchUsers"
        />
      </div>
      <ul>
        <li
          v-for="user in filteredUsers"
          :key="user.id"
          :class="{ active: activeChat === user.id }"
          @click="selectChat(user)"
          >
          {{ user.name }}
        </li>
      </ul>
    </div>

    <!-- 右侧聊天内容展示 -->
    <div class="chat-window">
      <div class="messages">
        <div
          v-for="(message, index) in activeMessages"
          :key="index"
          class="message"
        >
          <span>{{ message.sender }}: {{ message.content }}</span>
        </div>
      </div>
      <div class="input-area">
        <input
          type="text"
          v-model="newMessage"
          @keyup.enter="sendMessage"
          placeholder="Type your message here..."
        />
        <button @click="sendMessage">Send</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "ChatApp",
  data() {
    return {
      userName:"",
      searchQuery: "", // 搜索框输入内容
      users: [
        { id: "user1", name: "Alice" },
        { id: "user2", name: "Bob" },
        { id: "user3", name: "Charlie" },
      ], // 聊天对象列表
      filteredUsers: [], // 搜索过滤后的用户列表
      activeChat: null, // 当前选中的聊天对象 ID
      chatUserName: null, // 当前选中的聊天对象名称
      messages: {
        user1: [
          { sender: "Alice", content: "Hi there!" },
          { sender: "Me", content: "Hello!" },
        ],
        user2: [
          { sender: "Bob", content: "How are you?" },
          { sender: "Me", content: "I'm good, thanks!" },
        ],
        user3: [],
      }, // 所有聊天记录
      newMessage: "", // 新消息内容
    };
  },
  computed: {
    activeMessages() {
      // 当前选中聊天对象的消息
      return this.messages[this.activeChat] || [];
    },
  },
  methods: {
    searchUsers() {
      // 根据搜索框内容过滤用户列表
      this.filteredUsers = this.users.filter((user) =>
        user.name.toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    },
    selectChat(user) {
      // 选择聊天对象
      this.chatUserName = user.name;
      this.activeChat = user.id;
    },
    async sendMessage() {
      console.log("chatUserName", this.chatUserName)
      if (this.newMessage.trim() !== "") {
        // 将消息添加到本地消息列表
        this.messages[this.activeChat].push({
          sender: "Me",
          content: this.newMessage.trim(),
        });
        try {
          // 发送 HTTP 请求
          const response = await axios(
          {
            method: 'post',
            url: '/api/message/send',
            headers: {
              'Content-Type': 'application/json',
              'User-Name': sessionStorage.getItem('user_name')
            },
            data: {
              user_name: this.userName,
              receiver: this.chatUserName,
              content: this.newMessage,
              //content_type: this.contentType
            },
            withCredentials: true
          });
          console.log('Message sent successfully:', response.data);
          this.messages.push({
            sender: this.userName,
            receiver: this.chatUserName,
            content: this.newMessage.trim(),
            //content_type: this.contentType
          });
          this.newMessage = ''
        } catch (error) {
          console.error('Error sending message:', error);
        }

        // 清空输入框
        this.newMessage = "";
      }
    },
  },
  mounted() {
    //websocket
    this.userName = sessionStorage.getItem('user_name'); // 从 sessionStorage 获取 userID
    //TODO:获取不到怎么办？直接跳转回登录？
    this.socket = new WebSocket(`ws://localhost:8103/ws?userName=${this.userName}`); // 在 URL 中添加 userID

    // 监听 WebSocket 事件
    this.socket.onopen = () => {
      console.log('WebSocket connection established');
      if (this.socket.readyState !== WebSocket.OPEN) {
        setTimeout(() => {
          this.socket = new WebSocket('ws://localhost:8103/ws');
        }, 3000); // 3秒后重连
      } 
    };
    this.socket.onmessage = (event) => {
      console.log('Message received:', event.data);
      const message = JSON.parse(event.data);
      //TODO:这里面的sender其实应该是conversationID，暂时用sender区分不同的conversation
      
      // 检查 message.sender 是否存在于 this.messages 中
      if (!this.messages[message.sender]) {
        // 如果不存在，为该 sender 初始化一个空数组
        this.messages[message.sender] = [];
      }
      // 将接收到的消息添加到对应的聊天记录中
      this.messages[message.sender].push(message);
    };
    this.socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
    this.socket.onclose = () => {
      console.log('WebSocket connection closed, attempting to reconnect...');
      setTimeout(() => {
        this.socket = new WebSocket('ws://localhost:8103/ws');
      }, 3000); // 3秒后重连
    };



    // 默认显示第一个用户的聊天内容
    if (this.users.length > 0) {
      this.activeChat = this.users[0].id;
    }
    this.filteredUsers = this.users; // 初始化过滤后的用户列表
  },
};
</script>

<style>
.chat-container {
  display: flex;
  height: 100vh;
}

.chat-list {
  width: 30%;
  border-right: 1px solid #ccc;
  padding: 10px;
  background-color: #f9f9f9;
}

.search-bar {
  margin-bottom: 10px;
}

.search-bar input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.chat-list ul {
  list-style: none;
  padding: 0;
}

.chat-list li {
  padding: 10px;
  cursor: pointer;
  border-radius: 4px;
}

.chat-list li.active {
  background-color: #42b983;
  color: white;
}

.chat-list li:hover {
  background-color: #e1f5fe;
}

.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 10px;
}

.messages {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 10px;
}

.message {
  margin-bottom: 10px;
  padding: 8px;
  background-color: #e1f5fe;
  border-radius: 4px;
}

.input-area {
  display: flex;
}

.input-area input {
  flex: 1;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  margin-right: 10px;
}

.input-area button {
  padding: 8px 16px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.input-area button:hover {
  background-color: #369f6e;
}
</style>
















<!-- 

<template>
  <div class="chat-window">
    <div class="messages">
      <div v-for="(message, index) in messages" :key="index" class="message">
        <span>{{ message }}</span>
      </div>
    </div>
    <div class="input-area">
      <input
        type="text"
        v-model="newMessage"
        @keyup.enter="sendMessage"
        placeholder="Type your message here..."
      />
      <button @click="sendMessage">Send</button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
export default {
  name: 'ChatWindow',
  data() {
    return {
      newMessage: '',
      receiver: 'user123', // 接收者（默认值，可动态设置）
      contentType: 'text', // 消息类型（默认文字）
      socket: null, //websocket
      messages: [] //存储接收到和发送的消息
    }
  },
  mounted(){
    const userName = sessionStorage.getItem('user_name'); // 从 sessionStorage 获取 userID
    //TODO:获取不到怎么办？直接跳转回登录？
    this.socket = new WebSocket(`ws://localhost:8084/ws?userName=${userName}`); // 在 URL 中添加 userID

    // 监听 WebSocket 事件
    this.socket.onopen = () => {
      console.log('WebSocket connection established');
      if (this.socket.readyState !== WebSocket.OPEN) {
        setTimeout(() => {
          this.socket = new WebSocket('ws://localhost:8084/ws');
        }, 3000); // 3秒后重连
      } 
    };

    this.socket.onmessage = (event) => {
      console.log('Message received:', event.data);
      const message = JSON.parse(event.data);
      this.messages.push(message); // 将接收到的消息添加到消息列表
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    this.socket.onclose = () => {
      console.log('WebSocket connection closed, attempting to reconnect...');
      setTimeout(() => {
        this.socket = new WebSocket('ws://localhost:8084/ws');
      }, 3000); // 3秒后重连
    };
  },
  methods: {
    async sendMessage() {
      try {
          // 发送 HTTP 请求
          const response = await axios(
          {
            method: 'post',
            url: '/api/message/send',
            headers: {
              'Content-Type': 'application/json',
            },
            data: {
              user_name: this.username,
              receiver: this.receiver,
              content: this.newMessage,
              content_type: this.contentType
            },
            withCredentials: true
          });
          console.log('Message sent successfully:', response.data);
          this.messages.push({
            sender: sessionStorage.getItem('user_name'),
            receiver: this.receiver,
            content: this.newMessage.trim(),
            content_type: this.contentType
          });
          this.newMessage = ''
        } catch (error) {
          console.error('Error sending message:', error);
        }

      if (this.newMessage.trim() !== '') {
        this.messages.push(this.newMessage.trim())
        this.newMessage = ''
      }
      // 清空输入框
      this.newMessage = '';
    }
  }
}
</script>

<style>
.chat-window {
  display: flex;
  flex-direction: column;
  height: 100vh;
  max-width: 600px;
  margin: 0 auto;
  border: 1px solid #ccc;
  border-radius: 8px;
  overflow: hidden;
}

.messages {
  flex: 1;
  padding: 10px;
  overflow-y: auto;
  background-color: #f9f9f9;
}

.message {
  margin-bottom: 10px;
  padding: 8px;
  background-color: #e1f5fe;
  border-radius: 4px;
}

.input-area {
  display: flex;
  padding: 10px;
  border-top: 1px solid #ccc;
  background-color: #fff;
}

input {
  flex: 1;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  margin-right: 10px;
}

button {
  padding: 8px 16px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #369f6e;
}
</style> -->