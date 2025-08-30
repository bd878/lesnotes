import api from './api';
import models from './models';

async function moveMessage(id: number, thread: number) {
	let result = {
		error:  models.error(),
	}

	const form = new FormData()

	if (thread)
		form.append("thread", `${thread}`);

	if (id)
		form.append("id", `${id}`);

	try {
		const [_1, error] = await api("/messages/v1/update", {
			method: "POST",
			credentials: "include",
			body: form,
		});

		if (error)
			result.error = models.error(error)
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default moveMessage