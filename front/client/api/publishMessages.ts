import api from './api';
import models from './models';

async function publishMessages(ids: number[] = []) {
	let result = {
		error:  models.error(),
		ids:    [],
	}

	try {
		const [response, error] = await api("/messages/v1/publish", {
			queryParams: {
				ids: JSON.stringify(ids),
			},
			method: "PUT",
			credentials: "include",
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.ids = response.ids
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default publishMessages