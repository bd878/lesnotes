import api from './api';
import models from './models';

async function publishFileJson(token: string, id: number) {
	let result = {
		error:  models.error(),
	}

	try {
		const [_1, error] = await api("/files/v2/publish", {
			method: "POST",
			body: {
				token: token,
				req: {
					id: id,
				}
			}
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

export default publishFileJson