import api from './api';
import models from './models';

async function getMe() {
	let result = {
		error:  models.error(),
		user:   models.user(),
	}

	try {
		const [response, error] = await api("/users/v1/me", {
			credentials: 'include',
		});

		if (error)
			result.error = models.error(error)

		if (response)
			result.user = models.user(response)
	} catch(e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default getMe;
