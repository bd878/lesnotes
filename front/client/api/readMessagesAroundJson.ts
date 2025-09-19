import api from './api';
import models from './models';

async function readMessagesAroundJson(token: string, thread: number, id: number, limit: number) {
	let result = {
		error:       models.error(),
		messages:    [],
		isLastPage:  true,
		isFirstPage: true,
		total:       0,
		count:       0,
		offset:      0,
	}

	try {
		const [response, error] = await api('/messages/v2/read', {
			method: "POST",
			body: {
				token: token,
				req: {
					thread: thread,
					id:     id,
					limit:  limit,
				},
			},
		});

		if (error)
			result.error = models.error(error)

		if (response) {
			result.messages    = response.messages.map(models.message)
			result.isLastPage  = response.is_last_page
			result.isFirstPage = response.is_first_page
			result.total       = response.total
			result.count       = response.count
			result.offset      = response.offset
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readMessagesAroundJson;
