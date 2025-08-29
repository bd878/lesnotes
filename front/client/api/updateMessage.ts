import i18n from '../i18n';
import api from './api';
import * as is from '../third_party/is'
import models from './models';

async function updateMessage({id, text, public: isPublic}) {
	let response: any = {};
	let result: any = {
		error: "",
		explain: "",
		ID: "",
		updateUTCNano: "",
	}

	const form = new FormData()
	if (is.notEmpty(text))
		form.append("text", text);

	if (is.notUndef(isPublic))
		form.append("public", isPublic)

	try {
		response = await api("/messages/v1/update", {
			queryParams: {
				id: id,
			},
			method: "POST",
			credentials: "include",
			body: form,
		});

		if (response.error != "") {
			result.error = response.error
			result.explain = response.explain
		} else {
			if (response.value) {
				const model = models.message({update_utc_nano: response.value.update_utc_nano})
				result.ID = id
				result.updateUTCNano = model.updateUTCNano
			}
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = e
	}

	return result
}

export default updateMessage