import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ProfileView from '../views/ProfileView.vue'
import LoginView from '../views/LoginView.vue'
import SearchView from '../views/SearchView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/login', component: LoginView},
		{path: '/search', component: SearchView},
		{path: '/profile/:user_id', component: ProfileView},
	]
})

export default router
