import api from './api';
import models from './models';

async function deleteMessages(ids: number[] = []) {
	let result = {
		error:  models.error(),
		ids:    [],
	}

	try {
		const [response, error] = await api("/messages/v1/delete", {
			queryParams: {
				ids: JSON.stringify(ids),
			},
			method: "DELETE",
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

export default deleteMessages