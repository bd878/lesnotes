import api from './api';
import models from './models';

const EmptyListComments = {
	error:        models.error(),
	comments:     [],
}

async function listCommentsJson(token: string, message: number, name: string, limit: number, offset: number) {
	let result = {
		error:         models.error(),
		comments:      [],
	}

	console.log("listCommentsJson", "token", token, "message", message, "name", name, "limit", limit, "offset", offset)

	try {
		const [response, error] = await api('/comments/v2/list', {
			method: "POST",
			body: {
				token: token,
				req:   {
					message:  message,
					name:     name,
					limit:    limit,
					offset:   offset,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.comments = response.comments.map(models.comment)
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default listCommentsJson;
export { EmptyListComments };
