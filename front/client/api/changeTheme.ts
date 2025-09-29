import api from './api';
import models from './models';

async function changeTheme(theme: string) {
	let result = {
		error:     models.error(),
	}

	const form = new FormData()

	if (theme)
		form.append("theme", theme)

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

export default changeTheme;
