import i18n from '../i18n';
import api from './api';
import sendLog from './sendLog';
import models from './models';

async function readMessageJson(token, user, id) {
	let response = {};
	let result = {
		error:   false,
		explain: "",
		message: {},
	}

	try {
		await sendLog(token + " : " + JSON.stringify(req))

		response = await api('/messages/v2/read', {
			method: "POST",
			body: {
				token: token,
				req:   {
					user: user,
					id:   id,
				},
			},
		});

		if (response.error) {
			console.error('[readMessageJson]: /read response returned error')
			result.error = true
			result.explain = response.explain
		} else {
			result.message = models.message(response.messages[0])
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = e
	}

	return result;
}

export default readMessageJson;
