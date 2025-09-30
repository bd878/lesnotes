import api from './api';
import models from './models';

async function login(login: string, password: string) {
	let result = {
		error:   models.error(),
	}

	const form = new FormData()

	if (login)
		form.append("login", login);

	if (password)
		form.append("password", password);

	try {
		const [_1, error] = await api("/users/v1/login", {
			method: "POST",
			body:   form,
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

export default login

