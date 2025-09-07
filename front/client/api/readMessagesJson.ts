import api from './api';
import models from './models';

async function readMessagesJson(token: string, thread: number, order: number, limit: number, offset: number) {
	let result = {
		error:       models.error(),
		messages:    [],
		isLastPage:  false,
	}

	try {
		const [response, error] = await api('/messages/v2/read', {
			method: "POST",
			body: {
				token: token,
				req:   {
					thread: thread,
					limit:  limit,
					offset: offset,
					order:  order,
				},
			},
		});

		if (error)
			result.error = models.error(error)

		if (response) {
			result.messages = response.messages.map(models.message)
			result.isLastPage = response.isLastPage
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readMessagesJson;
