import { createStore } from 'vuex';

export default createStore({
  state: {
    userName: ''
  },
  mutations: {
    setUserName(state, userName) {
      state.userName = userName;
    }
  },
  actions: {
    login({ commit }, userName) {
      commit('setUserName', userName);
    }
  },
  getters: {
    userName: state => state.userName
  }
});