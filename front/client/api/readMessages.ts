import api from './api';
import models from './models';

async function readMessages(thread: number, order: number, limit: number, offset: number) {
	let result = {
		error:       models.error(),
		messages:    [],
		isLastPage:  true,
		isFirstPage: true,
	}

	try {
		const [response, error] = await api('/messages/v1/read', {
			queryParams: {
				thread: thread,
				limit:  limit,
				order:  order,
				offset: offset,
			},
			method: "GET",
			credentials: 'include',
		});

		if (error)
			result.error = models.error(error)

		if (response) {
			result.messages = response.messages.map(models.message)
			result.isLastPage = response.is_last_page
			if (offset > 0)
				result.isFirstPage = false
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readMessages;
