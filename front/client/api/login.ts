import api from './api';
import models from './models';

async function login(login: string, password: string, lang?: string) {
	let result = {
		error:     models.error(),
		token:     "",
		expiresAt: 0,
	}

	const form = new FormData()
	const headers = {}

	if (login) {
		form.append("login", login);
	}

	if (password) {
		form.append("password", password);
	}

	if (lang) {
		headers["X-Language"] = lang;
	}

	try {
		const [response, error] = await api("/users/v1/login", {
			method:  "POST",
			body:    form,
			headers: headers,
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.token = response.token
			result.expiresAt = response.expires_utc_nano
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default login

