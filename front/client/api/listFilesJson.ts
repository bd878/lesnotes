import api from './api';
import models from './models';

async function listFilesJson(token: string, user: number, limit: number, offset: number) {
	let result = {
		error:       models.error(),
		files:       [],
		paging:      models.paging(),
	}

	try {
		const [response, error] = await api('/files/v2/list', {
			method: "POST",
			body: {
				token: token,
				req:   {
					user:   user,
					limit:  limit,
					offset: offset,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.files  = response.files.map(models.file)
			result.paging = models.paging(response)
		}
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default listFilesJson;
