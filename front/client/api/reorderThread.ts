import api from './api';
import models from './models';

async function reorderThread(id: number, parent: number, next: number, prev: number) {
	let result = {
		error:  models.error(),
		id:     0,
	}

	const form = new FormData()

	if (id)
		form.append("id", id)

	if (parent)
		form.append("parent", parent)

	if (next)
		form.append("next", next)

	if (prev)
		form.append("prev", prev)

	try {
		const [response, error] = await api("/threads/v1/reorder", {
			method: "POST",
			credentials: "include",
			body: form,
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

export default reorderThread