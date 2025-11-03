import api from './api';
import models from './models';

async function changeFontSize(fontSize?: string) {
	let result = {
		error:     models.error(),
	}

	const form = new FormData()

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

	if (size)
		form.append("font_size", `${size}`)

	try {
		const [response, error] = await api("/users/v1/update", {
			method:      'POST',
			credentials: "include",
			body:        form,
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

export default changeFontSize;
