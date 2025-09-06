import api from './api';
import models from './models';

async function readMessagePathJson(token: string, id: number) {
	let result = {
		error:       models.error(),
		path:        [],
	}

	try {
		const [response, error] = await api('/messages/v2/read_path', {
			method: "POST",
			body: {
				token: token,
				req:   {
					id: id,
				},
			},
		});

		if (error)
			result.error = models.error(error)

		if (response)
			result.path = response.path.map(models.message)
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readMessagePathJson;
