import api from './api';
import models from './models';

async function logout() {
	let result = {
		error:  models.error(),
	}

	try {
		const [_1, error] = await api("/users/v1/logout", {
			method: 'POST',
			credentials: 'include',
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

export default logout;
