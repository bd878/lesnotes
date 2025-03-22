import {takeLatest,put,call,select} from 'redux-saga/effects'
import {
	UPDATE_MESSAGE,
	FETCH_MESSAGES,
	SEND_MESSAGE,
	DELETE_MESSAGE,
} from './messagesActions'
import {
	messagesFailedActionCreator,
	fetchMessagesSucceededActionCreator,
	sendMessageSucceededActionCreator,
	updateMessageSucceededActionCreator,
	deleteMessageSucceededActionCreator,
} from './messagesActionCreators'
import {selectMessages} from './messagesSelectors';
import api from '../../api'

interface FetchMessagesPayload {
	limit:  number;
	offset: number;
	order:  number;
}

function* fetchMessages({payload}: {payload: FetchMessagesPayload}) {
	try {
		const response = yield call(api.loadMessages,
			{limit: payload.limit, offset: payload.offset, order: payload.order})

		response.messages.reverse();
		if (response.error != "")
			yield put(messagesFailedActionCreator(response.error))
		else
			yield put(fetchMessagesSucceededActionCreator(response))
	} catch (e) {
		yield put(messagesFailedActionCreator(e.message))
	}
}

interface SendMessagePayload {
	message: any;
	file:    any;
}

function* sendMessage({payload}: {payload: SendMessagePayload}) {
	try {
		const response = yield call(api.sendMessage,
				{text: payload.message, file: payload.file})

		if (response.error != "")
			yield put(messagesFailedActionCreator(response.error))
		else
			yield put(sendMessageSucceededActionCreator(response.message))
	} catch (e) {
		yield put(messagesFailedActionCreator(e.message))
	}
}

function* updateMessage({payload}) {
	try {
		const response = yield call(api.updateMessage,
			payload.ID, payload.text)

		const messages = yield select(selectMessages)
		let idx = messages.findIndex(({ID}) => ID === payload.ID)
		if (idx !== -1)
			messages[idx] = {
				...messages[idx],
				ID: response.ID,
				text: payload.text,
				updateUTCNano: response.updateUTCNano,
			}

		if (response.error !== "")
			yield put(messagesFailedActionCreator(response.error))
		else
			yield put(updateMessageSucceededActionCreator(messages))
	} catch (e) {
		yield put(messagesFailedActionCreator(e.message))
	}
}

function* deleteMessage({payload}) {
	try {
		const response = yield call(api.deleteMessage, payload.ID)

		let messages = yield select(selectMessages)
		messages = messages.filter(({ID}) => ID !== payload.ID)

		if (response.error !== "")
			yield put(messagesFailedActionCreator(response.error))
		else
			yield put(deleteMessageSucceededActionCreator(messages))
	} catch (e) {
		yield put(messagesFailedActionCreator(e.message))
	}
}

function* messagesSaga() {
	yield takeLatest(DELETE_MESSAGE, deleteMessage)
	yield takeLatest(UPDATE_MESSAGE, updateMessage)
	yield takeLatest(FETCH_MESSAGES, fetchMessages)
	yield takeLatest(SEND_MESSAGE, sendMessage)
}

export {messagesSaga}