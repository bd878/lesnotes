import api from './api';
import models from './models';

async function changeLanguage(langCode: string) {
	let result = {
		error:     models.error(),
	}

	const form = new FormData()

	if (langCode)
		form.append("language", langCode)

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

export default changeLanguage;
