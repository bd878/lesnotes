import i18n from '../i18n';
import api from './api';
import * as is from '../third_party/is'
import models from './models';

async function moveMessage(id, threadID) {
	let response: any = {};
	let result = {
		error: "",
		explain: "",
		ID: "",
		threadID: 0,
		updateUTCNano: 0,
	}

	const form = new FormData()
	if (is.notUndef(threadID))
		form.append("thread_id", threadID);

	try {
		response = await api("/messages/v1/update", {
			queryParams: {
				id: id,
			},
			method: "POST",
			credentials: "include",
			body: form,
		});

		if (response.error !== "") {
			result.error = response.error
			result.explain = response.explain
		} else {
			if (response.value) {
				const model = models.message({update_utc_nano: response.value.update_utc_nano})
				result.ID = id
				result.threadID = threadID
				result.updateUTCNano = model.updateUTCNano
			}
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = e
	}

	return result
}

export default moveMessage