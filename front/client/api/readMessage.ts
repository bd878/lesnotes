import i18n from '../i18n';
import api from './api';
import models from './models';

async function readMessage(messageID: number) {
	let response: any = {};
	let result = {
		error: "",
		explain: "",
		message: {},
	}

	try {
		response = await api('/messages/v1/read', {
			queryParams: {
				id: messageID,
			},
			method: "GET",
			credentials: 'include',
		});

		if (response.error != "") {
			console.error('[readMessage]: /read response returned error')
			result.error = response.error
			result.explain = response.explain
		} else {
			result.message = models.message(response.value.message)
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = e
	}

	return result;
}

export default readMessage;
