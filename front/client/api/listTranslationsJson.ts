import api from './api';
import models from './models';

const EmptyListTranslations = {
	error:        models.error(),
	translations: models.translationPreview(),
}

async function listTranslationsJson(token: string, message: number, name?: string) {
	let result = {
		error:         models.error(),
		translations:  models.translationPreview(),
	}

	try {
		const [response, error] = await api('/translations/v2/list', {
			method: "POST",
			body: {
				token: token,
				req:   {
					message:  message,
					name:     name,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.translations = response.translations.map(models.translationPreview)
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default listTranslationsJson;
export { EmptyListTranslations };
