import i18n from '../i18n';
import api from './api';

async function validateMiniappData(body) {
	let result: SendMessageResult = {
		error: "",
		explain: "",
		ok: false,
		token: "",
	}

	try {
		let response = await api(`${BOT_VALIDATE_URL}`, {
			method: "POST",
			isFullUrl: true,
			body: body,
		});

		if (response.error != "") {
			result.error = response.error
			result.explain = response.explain
			result.ok = false
		} else {
			result.ok = true
			result.token = response.value.token
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = e.toString()
	}

	return result
}

export default validateMiniappData