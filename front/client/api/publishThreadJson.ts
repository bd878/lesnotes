import api from './api';
import models from './models';

async function publishThreadJson(token: string, id: number) {
	let result = {
		error:  models.error(),
		id:     0,
	}

	try {
		const [response, error] = await api("/threads/v2/publish", {
			method: "PUT",
			body:   {
				token:  token,
				req:    {
					id: id,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.id = response.id
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default publishThreadJson
