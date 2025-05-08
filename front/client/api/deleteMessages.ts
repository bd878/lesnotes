import i18n from '../i18n';
import api from './api';

async function deleteMessages(ids = []) {
	let response = {};
	let result: DeleteMessagesResult = {
		error: "",
		explain: "",
		IDs: [],
	}

	try {
		response = await api("/messages/v1/delete", {
			queryParams: {
				ids: ids,
			},
			method: "DELETE",
			credentials: "include",
		});

		if (response.error != "") {
			result.error = response.error
			result.explain = response.explain
		} else {
			if (response.value)
				result.IDs = ids
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		result.error = e
	}

	return result
}

export default deleteMessages