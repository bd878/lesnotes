import api from './api';
import models from './models';

async function searchMessagesJson(token: string, substr: string) {
	let result = {
		error:    models.error(),
		messages: [],
	}

	try {
		const [response, error] = await api('/messages/v2/search', {
			method: "POST",
			body:   {
				token:  token,
				req:    {
					query: substr,
				},
			},
		});

		if (error)
			result.error = models.error(error)

		if (response)
			result.messages = response.list.map(models.searchMessage)
	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default searchMessagesJson;
