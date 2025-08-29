import i18n from '../i18n';
import api from './api';
import models from './models';

async function readMessages(params) {
	const {
		limit,
		offset,
		order,
		threadID,
	} = params

	let response = {};
	let result = {
		error: "",
		explain: "",
		messages: [],
		isLastPage: false,
	}

	const queryParams = {
		limit: limit,
		offset: offset,
		asc: order,
	}

	if (threadID)
		queryParams.thread_id = threadID

	try {
		response = await api('/messages/v1/read', {
			queryParams: queryParams,
			method: "GET",
			credentials: 'include',
		});

		if (response.error != "") {
			console.error('[readMessages]: /read response returned error')
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

export default readMessages;
