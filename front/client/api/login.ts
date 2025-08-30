import api from './api';
import models from './models';

async function login(name: string, password: string) {
	let result = {
		error:   models.error(),
	}

	try {
		const [_1, error] = await api("/users/v1/login", {
			method: "POST",
			headers: {
				'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
			},
			body: new URLSearchParams({
				'name': name,
				'password': password,
			})
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

