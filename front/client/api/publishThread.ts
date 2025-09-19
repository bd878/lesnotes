import api from './api';
import models from './models';

async function publishThread(id: number) {
	let result = {
		error:  models.error(),
		id:     0,
	}

	try {
		const [response, error] = await api("/threads/v1/publish", {
			queryParams: {
				id: id,
			},
			method: "PUT",
			credentials: "include",
		});

		if (error)
			result.error = models.error(error)

		if (response)
			result.id = response.id
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default publishThread