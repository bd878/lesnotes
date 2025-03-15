import {
	SET_THREAD_MESSAGE,
  FETCH_MESSAGES,
  MESSAGES_FAILED,
  FETCH_MESSAGES_SUCCEEDED,
} from './threadsActions'

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

export function messagesFailedActionCreator(payload) {
  return {
    type: MESSAGES_FAILED,
    payload,
  }
}

export function fetchMessagesSucceededActionCreator(payload) {
  return {
    type: FETCH_MESSAGES_SUCCEEDED,
    payload,
  }
}