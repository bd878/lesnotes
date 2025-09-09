import api from './api';
import * as is from '../third_party/is'
import models from './models';

async function updateMessage(id: number, text?: string, title?: string, name?: string, fileIDs?: number[], isPublic?: boolean) {
	let result = {
		error:   models.error(),
	}

	const form = new FormData()

	if (text)
		form.append("text", text);

	if (id)
		form.append("id", `${id}`);

	if (title)
		form.append("title", title);

	if (name)
		form.append("name", name);

	if (fileIDs)
		form.append("file_ids", JSON.stringify(fileIDs))

	if (is.notUndef(isPublic))
		form.append("public", `${isPublic}`)

	try {
		const [_1, error] = await api("/messages/v1/update", {
			method: "POST",
			credentials: "include",
			body: form,
		});

		if (error)
			result.error = models.error(error)
	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default updateMessage