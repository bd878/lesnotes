import api from './api';
import models from './models';

// id = 0 : path = []
async function readPathJson(token: string, id: number) {
	let result = {
		error:       models.error(),
		path:        [],
	}

	try {
		const [response, error] = await api('/messages/v2/read_path', {
			method: "POST",
			body: {
				token: token,
				req:   {
					id: id,
				},
			},
		});

		if (error)
			result.error = models.error(error)

		// do not .reverse() here; append NullThread and then reverse
		if (response)
			result.path = response.path.map(models.thread)
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readPathJson;
