export default function getCurrentSession() {
	if (localStorage.getItem('token') == null) return sessionStorage.getItem('token');
    return localStorage.getItem('token');
}