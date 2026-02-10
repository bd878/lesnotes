import api from './api';
import models from './models';

const EmptyReadTranslation = {
	error:       models.error(),
	translation: models.translation(),
}

async function readTranslationJson(token: string, message: number, lang: string, name?: string) {
	let result = {
		error:         models.error(),
		translation:   models.translation(),
	}

	try {
		const [response, error] = await api('/translations/v2/read', {
			method: "POST",
			body: {
				token: token,
				req:   {
					message:  message,
					name:     name,
					lang:     lang,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.translation = models.translation(response.translation)
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readTranslationJson;
export { EmptyReadTranslation };
