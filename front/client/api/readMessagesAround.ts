import api from './api';
import models from './models';

async function readMessagesAround(threadID: number, id: number, limit: number) {
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
				thread: threadID,
				id:     id,
				limit:  limit,
			},
			method: "GET",
			credentials: 'include',
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

export default readMessagesAround;
