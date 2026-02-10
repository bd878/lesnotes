import api from './api';
import models from './models';

async function updateTranslationJson(token: string, message: number, lang: string, title: string, text: string) {
	let result = {
		error: models.error(),
	}

	try {
		const [_1, error] = await api("/translations/v2/update", {
			method: "PUT",
			body: {
				token: token,
				req: {
					message: message,
					lang:    lang,
					title:   title,
					text:    text,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		}
	} catch(e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default updateTranslationJson;
