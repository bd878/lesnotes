import api from './api';
import models from './models';

async function validateMiniappData(body: any) {
	let result = {
		error:   models.error(),
		token:   "",
	}

	try {
		const [response, error] = await api(`${BOT_VALIDATE_URL}`, {
			method: "POST",
			isFullUrl: true,
			body: body,
		});

		if (error)
			result.error = models.error(error)

		if (response)
			result.token = response.token
	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default validateMiniappData