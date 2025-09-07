import api from './api';
import models from './models';

async function readMessage(id: number) {
	let result = {
		error:    models.error(),
		message:  models.message(),
	}

	try {
		const [response, error] = await api('/messages/v1/read', {
			queryParams: {
				id: id,
			},
			method: "GET",
			credentials: 'include',
		});

		if (error)
			result.error = models.error(error)

		if (response)
			result.message = models.message(response.messages[0])
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readMessage;
