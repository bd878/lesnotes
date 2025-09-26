import api from './api';
import models from './models';

async function changeLanguageJson(token: string, langCode: string) {
	let result = {
		error:     models.error(),
	}

	try {
		const [response, error] = await api("/users/v2/update", {
			method:      'POST',
			body: {
				token:    token,
				req: {
					language: langCode,
				},
			},
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

export default changeLanguageJson;
