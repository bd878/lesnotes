import i18n from '../i18n';
import api from './api';
import * as is from '../third_party/is'
import models from './models';

interface SendMessageParams {
	text: any;
	file: any;
	threadID?: number;
}

async function sendMessage(params: SendMessageParams) {
	const {
		text,
		file,
		threadID,
	} = params

	let response = {};
	let result: SendMessageResult = {
		error: "",
		explain: "",
		message: "",
	}

	const form = new FormData()
	form.append("text", text);
	if (file != null && file.name != "") {
		form.append('file', file, file.name);
	}

	const queryParams = {}
	if (threadID !== 0 && is.notUndef(threadID))
		queryParams.thread_id = threadID

	try {
		response = await api("/messages/v1/send", {
			queryParams: queryParams,
			method: "POST",
			credentials: "include",
			body: form,
		});

		if (response.error != "") {
			result.error = response.error
			result.explain = response.explain
		} else {
			if (response.value != undefined && response.value.message != undefined)
				result.message = models.message(response.value.message)
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = e
	}

	return result
}

export default sendMessage