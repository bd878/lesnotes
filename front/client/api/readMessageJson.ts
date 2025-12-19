import api from './api';
import models from './models';

async function readMessageJson(token: string, user: number, id: number, name?: string) {
	let result = {
		error:     models.error(),
		message:   models.message(),
	}

	try {
		const [response, error] = await api('/messages/v2/read', {
			method: "POST",
			body: {
				token: token,
				req:   {
					user: user,
					id:   id,
					name: name,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.message = models.message(response.messages[0])
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readMessageJson;
