import i18n from '../i18n';
import api from './api';
import models from './models';

interface LoadMessagesResult {
	error: string;
	explain: string;
	messages: any[];
	isLastPage: boolean;
}

interface LoadMessagesParams {
	limit: number;
	offset: number;
	order: number;
	threadID: number;
}

async function loadMessages(params: LoadMessagesParams): LoadMessagesResult {
	const {
		limit,
		offset,
		order,
		threadID,
	} = params

	let response = {};
	let result: LoadMessagesResult = {
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
			console.error('[loadMessages]: /read response returned error')
			result.error = response.error
			result.explain = response.explain
		} else {
			result.messages = response.value.messages.map(models.message)
			result.isLastPage = response.value.is_last_page
		}
	} catch (e) {
		console.error(i18n("error_occured"), e);
		throw e
	}

	return result;
}

export default loadMessages;
