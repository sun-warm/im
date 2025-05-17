
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
          v-for="conversation in filteredConversations"
          :key="conversation.id"
          :class="{ active: activeConversation === conversation.id }"
          @click="selectConversation(conversation)"
        >
          <div class="user-item">
            <!-- 头像 -->
            <img :src="conversation.avatar" alt="avatar" class="avatar" />
            <div class="user-info">
              <!-- 名称 -->
              <div class="user-name">{{ conversation.name }}</div>
              <!-- 最近一条消息 -->
              <div class="last-message">{{ conversation.lastMessage || "No messages yet" }}</div>
            </div>
            <!-- 未读消息数 -->
            <div v-if="conversation.unreadCount > 0" class="unread-count">
              {{ conversation.unreadCount }}
            </div>
          </div>
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
      userName: "", // 当前用户名称
      searchQuery: "", // 搜索框输入内容
      conversations: [
        {
          id: "conv1",
          type: "single", // 单聊
          name: "Alice", // 对方名称
          avatar: "https://example.com/avatar1.png", // 对方头像
          lastMessage: "Hi there!", // 最近一条消息
          unreadCount: 2, // 未读消息数
        },
        {
          id: "conv2",
          type: "group", // 群聊
          name: "Project Team", // 群名称
          avatar: "https://example.com/group_avatar.png", // 群头像
          lastMessage: "Meeting at 3 PM",
          unreadCount: 5,
        },
      ], // 会话列表
      filteredConversations: [], // 搜索过滤后的会话列表
      activeConversation: null, // 当前选中的会话 ID
      messages: {
        conv1: [
          { sender: "Alice", content: "Hi there!" },
          { sender: "Me", content: "Hello!" },
        ],
        conv2: [
          { sender: "Bob", content: "Meeting at 3 PM" },
          { sender: "Me", content: "Got it!" },
        ],
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
    searchConversations() {
      this.filteredConversations = this.conversations.filter((conversation) =>
        conversation.name.toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    },
    selectConversation(conversation) {
      this.activeConversation = conversation.id;

      // 清除未读消息数
      const selectedConversation = this.conversations.find((c) => c.id === conversation.id);
      if (selectedConversation) {
        selectedConversation.unreadCount = 0;
      }
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
              conversation_id: this.activeConversation,
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
      console.log("Message received:", event.data);
      const message = JSON.parse(event.data);

      // 检查 message.conversation_id 是否存在于 this.messages 中
      if (!this.messages[message.conversation_id]) {
        // 如果不存在，为该会话初始化一个空数组
        this.messages[message.conversation_id] = [];
      }

      // 将接收到的消息添加到对应的聊天记录中
      this.messages[message.conversation_id].push(message);

      // 更新会话的最近消息和未读数
      const conversation = this.conversations.find((c) => c.id === message.conversation_id);
      if (conversation) {
        conversation.lastMessage = message.content; // 更新最近消息
        if (this.activeConversation !== message.conversation_id) {
          conversation.unreadCount += 1; // 如果不是当前会话，增加未读数
        }
      }
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

    // 默认显示第一个会话的聊天内容
    if (this.conversations.length > 0) {
      this.activeConversation = this.conversations[0].id;
    }
    this.filteredConversations = this.conversations; // 初始化过滤后的会话列表
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

.user-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border-radius: 4px;
  cursor: pointer;
}

.user-item:hover {
  background-color: #e1f5fe;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
}

.user-info {
  flex: 1;
}

.user-name {
  font-weight: bold;
  font-size: 16px;
}

.last-message {
  font-size: 14px;
  color: #888;
}

.unread-count {
  background-color: #f44336;
  color: white;
  font-size: 12px;
  font-weight: bold;
  padding: 4px 8px;
  border-radius: 12px;
  text-align: center;
  min-width: 20px;
  height: 20px;
  line-height: 20px;
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