import api from './api';
import models from './models';

async function changeFontSize(fontSize: number) {
	let result = {
		error:     models.error(),
	}

	const form = new FormData()

	if (fontSize)
		form.append("font_size", `${fontSize}`)

	try {
		const [response, error] = await api("/users/v1/update", {
			method:      'POST',
			credentials: "include",
			body:        form,
		});

		if (error)
			result.error = models.error(error)
	} catch(e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default changeFontSize;
