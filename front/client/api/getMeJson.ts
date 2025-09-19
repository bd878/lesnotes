import api from './api';
import models from './models';

async function getMeJson(token: string) {
	let result = {
		error:  models.error(),
		user:   models.user(),
	}

	try {
		const [response, error] = await api("/users/v2/me", {
			method: "POST",
			body: {
				token: token,
			},
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

export default getMeJson;
