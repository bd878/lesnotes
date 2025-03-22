import i18n from '../i18n';
import api from './api';
import models from './models';

interface LoadOneMessageResult {
	error: string;
	explain: string;
	message: any;
}

async function loadOneMessage(messageID: number): LoadOneMessageResult {
	let response = {};
	let result: LoadOneMessageResult = {
		error: "",
		explain: "",
		message: {},
	}

	try {
		response = await api('/messages/v1/read', {
			queryParams: {
				message_id: messageID,
			},
			method: "GET",
			credentials: 'include',
		});

		if (response.error != "") {
			console.error('[loadOneMessage]: /read response returned error')
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

export default loadOneMessage;
