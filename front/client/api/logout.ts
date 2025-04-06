import i18n from '../i18n';
import api from './api';

async function logout() {
	let response = {};
	let result = {
		error: "",
		explain: "",
	}

	try {
		response = await api("/users/v1/logout", {
			method: 'POST',
			credentials: 'include',
		});


		if (response.error == "") {
			/*ok*/
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

export default logout;
