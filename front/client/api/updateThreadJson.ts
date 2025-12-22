import api from './api';
import * as is from '../third_party/is'
import models from './models';

async function updateThreadJson(token: string, id: number, description?: string, name?: string) {
	let result = {
		error:   models.error(),
	}

	try {
		const [_1, error] = await api("/threads/v2/update", {
			method: "PUT",
			body: {
				token: token,
				req: {
					id:           id,
					description:  description,
					name:         name,
				}
			},
		});

		if (error.error) {
			result.error = models.error(error)
		}
	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default updateThreadJson