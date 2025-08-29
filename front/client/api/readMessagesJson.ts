import i18n from '../i18n';
import api from './api';
import sendLog from './sendLog';
import models from './models';

async function readMessagesJson(token, req) {
	let response: any = {};
	let result = {
		error: "",
		explain: "",
		messages: [],
		isLastPage: false,
	}

	try {
		await sendLog(token + " : " + JSON.stringify(req))

		response = await api('/messages/v2/read', {
			method: "POST",
			body: {
				token: token,
				req: req,
			},
		});

		if (response.error != "") {
			console.error('[readMessagesJson]: /read response returned error')
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

export default readMessagesJson;
