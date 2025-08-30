import api from './api';
import models from './models';

async function privateMessages(ids: number[] = []) {
	let result = {
		error:  models.error(),
		ids:    [],
	}

	try {
		const [response, error] = await api("/messages/v1/private", {
			queryParams: {
				ids: JSON.stringify(ids),
			},
			method: "PUT",
			credentials: "include",
		});

		if (error)
			result.error = models.error(error)

		if (response)
			result.ids = response.ids
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default privateMessages
