import api from './api';
import models from './models';

async function readMessagesJson(token: string, user: number, thread: number, order: number, limit: number, offset: number) {
	let result = {
		error:       models.error(),
		messages:    [],
		paging:      models.paging(),
	}

	try {
		const [response, error] = await api('/messages/v2/read', {
			method: "POST",
			body: {
				token: token,
				req:   {
					thread: thread,
					user:   user,
					limit:  limit,
					offset: offset,
					order:  order,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.messages = response.messages.map(models.message)
			result.paging   = models.paging(response)
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readMessagesJson;
