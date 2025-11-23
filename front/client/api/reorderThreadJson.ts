import api from './api';
import models from './models';

async function reorderThreadJson(token: number, id: number, parent: number, next: number, prev: number) {
	let result = {
		error:  models.error(),
		id:     0,
	}

	try {
		const [response, error] = await api("/threads/v2/reorder", {
			method: "POST",
			body: {
				token: token,
				req: {
					id:     id,
					parent: parent,
					next:   next,
					prev:   prev,
				},
			},
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

export default reorderThreadJson