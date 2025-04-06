import i18n from '../i18n';
import api from './api';

async function login(name, password) {
	let response = {}
	let result = {
		error: "",
		explain: "",
		isOk: false,
	}

	try {
		response = await api("/users/v1/login", {
			method: "POST",
			headers: {
				'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
			},
			body: new URLSearchParams({
				'name': name,
				'password': password,
			})
		});

		if (response.error == "") {
			if (response.value.status == "ok") {
				result.isOk = true
			} else {
				result.isOk = false
				result.explain = result.value.status
			}
		} else {
			console.error(i18n("error_occured"), response.error, response.explain)
			result.error = response.error
			result.explain = response.explain
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = e
	}

	return result
}

export default login

