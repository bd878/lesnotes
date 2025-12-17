import api from './api';
import models from './models';

async function sendMessageJson(token: string, text: string, title: string, fileIDs: number[], thread: number, isPrivate: boolean) {
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
					text:       text,
					title:      title,
					thread:     thread,
					file_ids:   fileIDs,
					private:    isPrivate,
				},
			},
		});

		if (error.error)
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
