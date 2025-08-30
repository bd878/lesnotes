import api from './api';
import models from './models';

async function authJson(token: string) {
	let result = {
		error:     models.error(),
		expired:   true,
	}

	try {
		const [response, error] = await api("/users/v2/auth", {
			method:      'POST',
			body:  {
				token: token,
			}
		});

		if (error)
			result.error = models.error(error)

		if (response)
			result.expired = response.expired
	} catch(e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default authJson;
