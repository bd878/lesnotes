import api from './api';
import models from './models';

async function changeFontSizeJson(token: string, fontSize: string) {
	let result = {
		error:     models.error(),
	}

	let size: number = 10
	switch (fontSize) {
	case "small":
		size = 8;
		break;
	case "medium":
		size = 10;
		break;
	case "large":
		size = 16;
		break;
	}

	try {
		const [response, error] = await api("/users/v2/update", {
			method:      'POST',
			body: {
				token:    token,
				req: {
					font_size: size,
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

export default changeFontSizeJson;
