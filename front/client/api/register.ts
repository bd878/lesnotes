import i18n from '../i18n';
import api from './api';

async function register(name, password) {
	let response: any = {}
	const result = {
		error: "",
		explain: "",
		isOk: false,
	}

	try {
		response = await api("/users/v1/signup", {
			method: "POST",
			headers: {
				'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
			},
			body: new URLSearchParams({
				'name': name,
				'password': password,
			})
		});

		if (response.error != "") {
			console.error("[RegisterForm]: /signup response returned error", response.error, response.explain)
			result.error = response.error
			result.explain = response.explain
		} else {
			if (response.value.status == "ok") {
				result.isOk = true
			} else {
				result.isOk = false
				result.explain = response.value.status
			}
		}
	} catch (e) {
		console.error(i18n("error_occured"), e)
		result.error = e
	}

	return result
}

export default register