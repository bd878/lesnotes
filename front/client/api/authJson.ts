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

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.expired = response.expired
		}
	} catch(e) {
		console.error(e)
		result.error.error   = true
		result.error.code    = 911
		result.error.status  = 500
	}

	return result
}

export default authJson;
