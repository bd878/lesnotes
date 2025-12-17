import api from './api';
import mapMessage from './models/message';
import mapError from './models/error';

async function sendMessage(text: string, title?: string, fileIDs?: number[], thread?: number) {
	let result = {
		error:   mapError(),
		message: mapMessage(),
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
			result.error   = mapError(error)

		if (response)
			result.message = mapMessage(response.message)
	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default sendMessage