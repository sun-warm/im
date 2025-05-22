
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
       <!-- 加号按钮 -->
       <div class="add-friend-icon" @click="toggleMenu">
          <i class="fas fa-plus"></i> <!-- 使用 Font Awesome 加号图标 -->
          <!-- 弹出菜单 -->
          <div v-if="showMenu" class="dropdown-menu">
            <ul>
              <li @click="openAddFriendModal">
                <i class="fas fa-user-plus"></i> 添加好友
              </li>
              <li @click="openCreateChatModal">
                <i class="fas fa-users"></i> 创建聊天
              </li>
            </ul>
          </div>
        </div>
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

    <!-- 添加好友子页面 -->
    <div v-if="showAddFriendModal" class="add-friend-modal">
      <div class="modal-content">
        <h3>添加好友</h3>
        <input
          type="text"
          v-model="addFriendQuery"
          placeholder="输入用户名搜索好友..."
          @input="searchUsers"
        />
        <ul>
          <li
            v-for="user in searchResults"
            :key="user.id"
            @click="addFriend(user)"
          >
            {{ user.name }}
          </li>
        </ul>
        <button @click="showAddFriendModal = true">添加好友</button>
        <button @click="showAddFriendModal = false">取消</button>
      </div>
    </div>

    <!-- 创建聊天窗口 --> 
    <div v-if="showCreateChatModal" class="create-chat-modal">
      <div class="modal-content">
        <h3>创建聊天</h3>
        <input
          type="text"
          v-model="createChatQuery"
          placeholder="搜索好友..."
          @input="filterFriends"
        />
        <ul>
          <li
            v-for="friend in filteredFriends"
            :key="friend.id"
          >
            <label>
              <input
                type="checkbox"
                :value="friend.id"
                v-model="selectedFriends"
              />
              {{ friend.name }}
            </label>
          </li>
        </ul>
        <button @click="createChat" class="btn">创建</button>
        <button @click="showCreateChatModal = false" class="btn-link">取消</button>
      </div>
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
      showMenu: false, // 控制弹出菜单的显示
      userName: "", // 当前用户名称
      searchQuery: "", // 搜索框输入内容

      showCreateChatModal: false, // 控制创建聊天窗口的显示
      createChatQuery: "", // 搜索好友的输入内容
      Friends: [], // 好友列表

      addFriendQuery: "", // 添加好友时的搜索框内容
      showAddFriendModal: false, // 控制添加好友子页面的显示
      searchResults: [], // 搜索好友的结果

      conversations: [
        {
          id: "conv1",
          type: "single", // 单聊showAddFriendModal
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
    toggleMenu() {
      this.showMenu = !this.showMenu; // 切换菜单显示状态
    },
    openAddFriendModal() {
      this.showMenu = false; // 关闭菜单
      this.showAddFriendModal = true; // 显示添加好友窗口
    },
    openCreateChatModal() {
      this.showMenu = false; // 关闭菜单
      this.showCreateChatModal = true; // 显示创建聊天窗口
    },
    filterFriends() {
      // 根据搜索内容过滤好友列表
      this.filteredFriends = this.friends.filter((friend) =>
        friend.name.toLowerCase().includes(this.createChatQuery.toLowerCase())
      );
    },
    createChat() {
      if (this.selectedFriends.length === 0) {
        alert("请选择至少一个好友！");
        return;
      }
      console.log("创建聊天，选中的好友：", this.selectedFriends);
      // 调用后端接口创建群聊
      // 示例：
      // axios.post('/api/create-chat', { members: this.selectedFriends })
      //   .then(response => {
      //     console.log('群聊创建成功:', response.data);
      //     this.showCreateChatModal = false; // 关闭窗口
      //   })
      //   .catch(error => {
      //     console.error('创建群聊失败:', error);
      //   });
      this.showCreateChatModal = false; // 关闭窗口
    },
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

    addFriend(user) {
      // 添加好友逻辑，可以调用后端接口
      console.log("Adding friend:", user);
      this.conversations.push({
        id: user.id,
        type: "single",
        name: user.name,
        avatar: "https://example.com/default_avatar.png",
        lastMessage: "",
        unreadCount: 0,
      });
      this.showAddFriendModal = false; // 关闭子页面
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
  box-sizing: border-box; /* 确保 padding 不会导致宽度超出 */
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
  display: flex; /* 使用 flex 布局 */
  align-items: center; /* 垂直居中 */
  margin-bottom: 10px;
  box-sizing: border-box; /* 确保 padding 不会导致宽度超出父容器 */
  width: 100%; /* 确保宽度不超出父容器 */
}

.search-bar input {
  flex: 1; /* 输入框占据剩余空间 */
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-sizing: border-box; /* 确保 padding 不会影响宽度 */
}

.add-friend-icon {
  margin-left: 10px;
  padding: 8px;
  background-color: transparent;
  border: none;
  cursor: pointer;
  font-size: 18px;
  color: #42b983;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.add-friend-icon:hover {
  color: #369f6e;
}

.add-friend-modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.dropdown-menu {
  position: absolute;
  top: 40px; /* 紧贴加号按钮的下方 */
  left: 0;
  /*right: 0;*/
  background-color: white;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  width: 150px;
  /*position: relative;*/
}

.dropdown-menu ul {
  list-style: none;
  margin: 0;
  padding: 0;
}

.dropdown-menu li {
  display: flex; /* 使用 flex 布局 */
  align-items: center; /* 垂直居中 */
  padding: 10px;
  cursor: pointer;
  text-align: left; /* 左对齐 */
}

.dropdown-menu li i {
  margin-right: 8px; /* 图标与文字之间的间距 */
  font-size: 16px; /* 图标大小 */
}

.dropdown-menu li:hover {
  background-color: #f5f5f5;
}

.modal-content {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.modal-content input {
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.modal-content ul {
  list-style: none;
  padding: 0;
}

.modal-content li {
  padding: 8px;
  cursor: pointer;
  border-radius: 4px;
}

.modal-content li:hover {
  background-color: #e1f5fe;
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
