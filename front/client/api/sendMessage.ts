import i18n from '../i18n';
import api from './api';
import * as is from '../third_party/is'
import models from './models';

async function sendMessage(params) {
	const {
		text,
		fileID,
		threadID,
	} = params

	let response: any = {};
	let result = {
		error: "",
		explain: "",
		message: {},
	}

	const form = new FormData()
	form.append("text", text);
	if (is.notUndef(fileID))
		form.append("file_id", fileID);

	const queryParams: any = {}
	if (is.notEmpty(threadID))
		queryParams.thread = threadID

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