import api from './api';
import models from './models';

async function deleteFileJson(token: string, id: number) {
	let result = {
		error:  models.error(),
	}

	try {
		const [_1, error] = await api("/files/v2/delete", {
			method: "DELETE",
			body: {
				token: token,
				req: {
					ids: [id],
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default deleteFileJson