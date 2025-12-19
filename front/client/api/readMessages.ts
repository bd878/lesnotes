import api from './api';
import models from './models';

async function readMessages(thread: number, order: number, limit: number, offset: number) {
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

		if (error.error) {
			result.error = models.error(error)
		} else {
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

export default readMessages;
