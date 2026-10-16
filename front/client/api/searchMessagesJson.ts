import api from './api';
import models from './models';

async function searchMessagesJson(token: string, query: string, thread: number, isPublic: number) {
	let result = {
		error: models.error(),
		list: models.EmptySearchList,
	}

	console.log("searchMessagesJson", "token", token, "query", query, "thread", thread, "public", isPublic)

	try {
		const [response, error] = await api("/search/v2/messages", {
			method: "POST",
			body: {
				token: token,
				req: {
					query: query,
					public: isPublic,
					thread: thread,
				}
			}
		})

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.list = models.searchMessagesList(response)
		}
	} catch (e) {
		console.log(e)
		result.error.error = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default searchMessagesJson;
