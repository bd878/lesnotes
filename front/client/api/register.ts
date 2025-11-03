import models from './models';
import api from './api';

async function register(login: string, password: string, lang?: string) {
	let result = {
		error:       models.error(),
	}

	const form = new FormData()
	const headers = {}

	if (login)
		form.append("login", login);

	if (password)
		form.append("password", password);

	if (lang)
		headers["X-Language"] = lang;

	try {
		const [_1, error] = await api("/users/v1/signup", {
			method:  "POST",
			body:    form,
			headers: headers,
		});

		if (error)
			result.error = models.error(error)
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default register