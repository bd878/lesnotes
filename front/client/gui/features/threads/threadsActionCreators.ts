import {
	RESET,
	FAILED,
	SEND_MESSAGE,
	SET_THREAD_MESSAGE,
	FETCH_MESSAGES,
	FETCH_MESSAGES_SUCCEEDED,
	SEND_MESSAGE_SUCCEEDED,
} from './threadsActions'

// export const __setThreadMessageActionCreator = payload => ({action:__SET_THREAD_MESSAGE, payload}) 

export function setThreadMessageActionCreator(payload) {
	return {
		type: SET_THREAD_MESSAGE,
		payload,
	}
}

export function fetchMessagesActionCreator(limit: number, offset: number, order: number) {
	return {
		type: FETCH_MESSAGES,
		payload: {limit, offset, order},
	}
}

export function fetchMessagesSucceededActionCreator(payload) {
	return {
		type: FETCH_MESSAGES_SUCCEEDED,
		payload,
	}
}

export function resetActionCreator(_payload) {
	return {
		type: RESET,
		payload: {},
	}
}

export function sendMessageActionCreator(message, file) {
	return {
		type: SEND_MESSAGE,
		payload: {message, file},
	}
}

export function failedActionCreator(payload) {
	return {
		type: FAILED,
		payload,
	}
}

export function sendMessageSucceededActionCreator(payload) {
	return {
		type: SEND_MESSAGE_SUCCEEDED,
		payload,
	}
}