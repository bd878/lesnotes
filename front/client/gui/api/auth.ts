import i18n from '../i18n';
import api from './api';
import models from './models';

async function auth() {
	let response = {};
	let result = {
		error: "",
		explain: "",
		expired: false,
	}

	try {
		response = await api("/users/v1/auth", {
			method: 'POST',
			credentials: 'include',
		});

		if (response.error == "") {
			if (response.value.expired)
				result.expired = true
			else
				result.user = models.user(response.value.user)
		} else {
			result.error = response.error
			result.explain = response.explain
		}

	} catch(e) {
		console.error(i18n("error_occured"), e);
		result.error = e
	}

	return result
}

export default auth;
