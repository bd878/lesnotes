import api from './api';
import models from './models';

async function auth() {
	let result = {
		error:     models.error(),
		expired:   true,
	}

	try {
		const [response, error] = await api("/users/v1/auth", {
			method:      'POST',
			credentials: 'include',
		});

		if (error)
			result.error = models.error(error)

		if (response)
			result.expired = response.expired
	} catch(e) {
		result.error.error   = true
		result.error.code    = 912
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default auth;
