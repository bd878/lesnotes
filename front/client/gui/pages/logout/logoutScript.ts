import api from '../../../api';

async function init() {
	await api.logout()
	setTimeout(() => { location.href = "/login" }, 0)
}

window.addEventListener("load", () => {
	init();
})