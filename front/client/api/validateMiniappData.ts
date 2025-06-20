import i18n from '../i18n';
import api from './api';
import * as is from '../third_party/is'
import models from './models';

async function validateMiniappData(body) {
	const {
		text,
		fileID,
		threadID,
	} = params

	let response = {};
	let result: SendMessageResult = {
		error: "",
		explain: "",
		ok: false,
	}

	try {
		response = await api(`${BOT_VALIDATE_URL}`, {
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
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = e
	}

	return result
}

export default validateMiniappData