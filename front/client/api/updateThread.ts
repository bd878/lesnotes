import api from './api';
import * as is from '../third_party/is'
import models from './models';

async function updateThread(id: number, description: string, title: string, name?: string) {
	let result = {
		error:   models.error(),
	}

	const form = new FormData()

	if (description) {
		form.append("description", description);
	}

	if (id) {
		form.append("id", `${id}`);
	}

	if (title) {
		form.append("title", title);
	}

	if (name) {
		form.append("name", name);
	}

	console.log("updateThread", "id", id, "description", description, "title", title, "name", name)

	try {
		const [_1, error] = await api("/threads/v1/update", {
			method: "PUT",
			credentials: "include",
			body: form,
		});

		if (error.error) {
			result.error = models.error(error)
		}
	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default updateThread