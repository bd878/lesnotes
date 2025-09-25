import api from './api';
import models from './models';

async function changeThemeJson(token: string, theme: string) {
	let result = {
		error:     models.error(),
	}

	try {
		const [response, error] = await api("/users/v2/update", {
			method:      'POST',
			body: {
				token:    token,
				req: {
					theme: theme,
				},
			},
		});

		if (error)
			result.error = models.error(error)
	} catch(e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default changeThemeJson;
