import api from './api';
import models from './models';

async function sendCommentJson(token: string, message: number, text: string) {
	let result = {
		error:    models.error(),
		id:       0,
	}

	console.log("sendCommentJson", "message_id", message, "text", text)

	try {
		const [response, error] = await api("/comments/v2/send", {
			method: "POST",
			body: {
				token: token,
				req: {
					text:    text,
					message: message,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.id = response.id
		}
	} catch (e) {
		result.error.error     = true
		result.error.status    = 500
		result.error.explain   = e.toString()
	}

	return result
}

export default sendCommentJson
