import api from './api';
import models from './models';

async function readMessagesAround(threadID: number, id: number, limit: number) {
	let result = {
		error:       models.error(),
		messages:    [],
		isLastPage:  false,
		isFirstPage: false,
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
			result.messages = response.messages.map(models.message)
			result.isLastPage = response.isLastPage
			result.isFirstPage = response.isFirstPage
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readMessagesAround;
