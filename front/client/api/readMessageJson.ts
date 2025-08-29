import i18n from '../i18n';
import api from './api';
import sendLog from './sendLog';
import models from './models';

async function readMessageJson(token, user, id) {
	let result = {
		error:   false,
		status:  200,
		explain: "",
		code:    0,
		message: models.message(),
	}

	try {
		const response = await api('/messages/v2/read', {
			method: "POST",
			body: {
				token: token,
				req:   {
					user: user,
					id:   id,
				},
			},
		});

		result.error    = response.error
		result.status   = response.status
		result.explain  = response.explain
		result.code     = response.code

		if (response.error) {
			console.error('[readMessageJson]: read response returned error', "status:", response.status, "code:", response.code, "explain:", response.explain)

			return result
		} else {
			result.message = models.message(response.messages[0])
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = 500
		result.explain = e
	}

	return result;
}

export default readMessageJson;
