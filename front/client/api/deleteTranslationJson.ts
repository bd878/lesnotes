import api from './api';
import models from './models';

async function deleteTranslationJson(token: string, message: number, lang: string) {
	let result = {
		error: models.error(),
	}

	console.log("deleteTranslationJson", "token", token, "message", message, "lang", lang)

	try {
		const [_1, error] = await api("/translations/v2/delete", {
			method: "DELETE",
			body: {
				token: token,
				req: {
					message: message,
					lang:    lang,
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

export default deleteTranslationJson;
