import api from './api';
import models from './models';

async function sendMessageJson(token: string, text: string, file: number, thread: number) {
	let result = {
		error:   models.error(),
		message: models.message(),
	}

	try {
		const [response, error] = await api("/messages/v2/send", {
			method: "POST",
			body: {
				token: token,
				req: {
					text:   text,
					file:   file,
					thread: thread,
				},
			},
		});

		if (error)
			result.error = models.error(error)

		if (response)
			result.message = models.message(response.message)

	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default sendMessageJson
