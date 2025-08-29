import i18n from '../i18n';
import api from './api';
import models from './models';

async function privateMessages(ids = []) {
	let response: any = {};
	let result: any = {
		error: "",
		explain: "",
		IDs: [],
	}

	try {
		response = await api("/messages/v1/private", {
			queryParams: {
				ids: JSON.stringify(ids),
			},
			method: "PUT",
			credentials: "include",
		});

		if (response.error != "") {
			result.error = response.error
			result.explain = response.explain
		} else {
			if (response.value) {
				const model = models.message({update_utc_nano: response.value.update_utc_nano})
				result.IDs = response.value.ids
				result.updateUTCNano = model.updateUTCNano
			}
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = e
	}

	return result
}

export default privateMessages