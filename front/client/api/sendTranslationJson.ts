import api from './api';
import models from './models';

async function sendTranslationJson(token: string, message: number, lang: string, text: string, title: string) {
	let result = {
		error:       models.error(),
		translation: models.translation(),
	}

	try {
		const [response, error] = await api("/translations/v2/send", {
			method: "POST",
			body: {
				token: token,
				req: {
					text:       text,
					title:      title,
					message:    message,
					lang:       lang,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.translation = models.translation(response.translation)
		}

	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default sendTranslationJson
