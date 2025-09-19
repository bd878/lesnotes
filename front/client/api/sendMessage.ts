import api from './api';
import models from './models';

async function sendMessage(text: string, title?: string, fileIDs?: number[], thread?: number) {
	let result = {
		error:   models.error(),
		message: models.message(),
	}

	const form = new FormData()

	if (text)
		form.append("text", text);

	if (fileIDs)
		form.append("file_ids", JSON.stringify(fileIDs));

	if (title)
		form.append("title", title);

	if (thread)
		form.append("thread", `${thread}`)

	try {
		const [response, error] = await api("/messages/v1/send", {
			method:      "POST",
			credentials: "include",
			body:        form,
		});

		if (error)
			result.error   = models.error(error)

		if (response)
			result.message = models.message(response.message)
	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default sendMessage