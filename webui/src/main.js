import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import { axios, updateToken as axiosUpdate } from './services/axios.js';
import getCurrentSession from './services/authentication';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import PostCard from './components/PostCard.vue'
import UserCard from './components/UserCard.vue'
import ProfileCounters from './components/ProfileCounters.vue'
import Modal from './components/Modal.vue'
import IntersectionObserver from './components/IntersectionObserver.vue'
// import 'bootstrap-icons/font/bootstrap-icons.css'

import './assets/dashboard.css'
import './assets/main.css'
import 'bootstrap-icons/font/bootstrap-icons.css'

// Create the Vue SPA
const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.config.globalProperties.$axiosUpdate = axiosUpdate;
app.config.globalProperties.$currentSession = getCurrentSession;

// Register the components
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("PostCard", PostCard);
app.component("UserCard", UserCard);
app.component("ProfileCounters", ProfileCounters);
app.component("Modal", Modal);
app.component("IntersectionObserver", IntersectionObserver);

app.use(router)
app.mount('#app')