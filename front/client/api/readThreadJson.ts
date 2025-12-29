import api from './api';
import models from './models';

async function readThreadJson(token: string, user: number, id: number, name?: string) {
	let result = {
		error:     models.error(),
		thread:    models.thread(),
	}

	try {
		const [response, error] = await api('/threads/v2/read', {
			method: "POST",
			body: {
				token: token,
				req:   {
					id:   id,
					name: name,
					user: user,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.thread = models.thread(response.thread)
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readThreadJson;
