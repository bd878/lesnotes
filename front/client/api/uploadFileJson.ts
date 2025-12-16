import api from './api';
import models from './models';

async function uploadFileJson(token: string, file: any) {
	let result = {
		error:   models.error(),
		ID:      0,
		name:    "",
	}

	const form = new FormData()
	if (file)
		form.append('file', file, file.name);

	try {
		const [response, error] = await api("/files/v2/upload", {
			method: "POST",
			body: form,
		});

		if (error.error)
			result.error = models.error(error)

		if (response) {
			result.ID = response.id
			result.name = response.name
		}
	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default uploadFileJson