import api from './api';
import models from './models';

async function readBatchMessagesJson(token: string, ids: number[]) {
	let result = {
		error:       models.error(),
		messages:    [],
	}

	try {
		const [response, error] = await api('/messages/v2/read', {
			method: "POST",
			body: {
				token: token,
				req:   {
					ids: ids,
				},
			},
		});

		if (error)
			result.error = models.error(error)

		if (response) {
			result.messages = response.messages.map(models.message)
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readBatchMessagesJson;
