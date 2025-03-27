import i18n from '../i18n';
import api from './api';
import models from './models';

async function uploadFile(file: any) {
	let response = {};
	let result: uploadFileResult = {
		error: "",
		explain: "",
		ID: "",
		updateUTCNano: 0,
	}

	const form = new FormData()
	if (file != null && file.name != "") {
		form.append('file', file, file.name);
	} else {
		result.error = "file required"
		return result
	}

	try {
		response = await api("/files/v1/upload", {
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

export default uploadFile