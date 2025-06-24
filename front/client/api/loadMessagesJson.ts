import i18n from '../i18n';
import api from './api';
import models from './models';

interface LoadMessagesResult {
	error: string;
	explain: string;
	messages: any[];
	isLastPage: boolean;
}

async function loadMessagesJson(token, req): LoadMessagesResult {
	let response = {};
	let result: LoadMessagesResult = {
		error: "",
		explain: "",
		messages: [],
		isLastPage: false,
	}

	api.sendLog(token + " : " + JSON.stringify(req))

	try {
		response = await api('/messages/v2/read', {
			method: "POST",
			body: {
				token: token,
				req: req,
			},
		});

		if (response.error != "") {
			console.error('[loadMessagesJson]: /read response returned error')
			result.error = response.error
			result.explain = response.explain
		} else {
			result.messages = response.value.messages.map(models.message)
			result.isLastPage = response.value.is_last_page
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = e
	}

	return result;
}

export default loadMessagesJson;
