import api from './api';
import models from './models';

const EmptyReadComment = {
	error:   models.error(),
	comment: models.comment(),
}

async function readCommentJson(token: string, id: number) {
	let result = {
		error:     models.error(),
		comment:   models.comment(),
	}

	console.log("readCommentJson", "token", token, "id", id)

	try {
		const [response, error] = await api('/comments/v2/read', {
			method: "POST",
			body: {
				token: token,
				req:   {
					id:   id,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.comment = models.comment(response.comment)
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readCommentJson;
export { EmptyReadComment };
