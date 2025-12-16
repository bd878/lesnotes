import api from './api';
import * as is from '../third_party/is'
import models from './models';

async function updateMessageJson(token: string, id: number, text?: string, title?: string, name?: string, fileIDs?: number[]) {
	let result = {
		error:   models.error(),
	}

	try {
		const [_1, error] = await api("/messages/v2/update", {
			method: "POST",
			body: {
				token: token,
				req: {
					id:        id,
					text:      text,
					title:     title,
					name:      name,
					file_ids:  fileIDs,
				}
			},
		});

		if (error.error)
			result.error = models.error(error)
	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default updateMessageJson