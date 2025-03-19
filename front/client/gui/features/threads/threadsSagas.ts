import {takeLatest,takeEvery,put,call,select} from 'redux-saga/effects'
import {FETCH_MESSAGES, SEND_MESSAGE, SET_THREAD_MESSAGE, __SET_THREAD_MESSAGE} from './threadsActions'
import {selectThreadID} from './threadsSelectors'
import {
	sendMessageActionCreator,
	failedActionCreator,
	fetchMessagesSucceededActionCreator,
	sendMessageSucceededActionCreator,
	failedActionCreator,
	resetActionCreator,
} from './threadsActionCreators'
import api from '../../api'

interface FetchMessagesPayload {
	limit:  number;
	offset: number;
	order:  number;
}

function* fetchMessages({payload}: {payload: FetchMessagesPayload}) {
	try {
		const threadID = yield select(selectThreadID)
		const response = yield call(api.loadMessages,
			{limit: payload.limit, offset: payload.offset, order: payload.order, threadID: threadID})

		response.messages.reverse();
		if (response.error != "")
			yield put(failedActionCreator(response.error))
		else
			yield put(fetchMessagesSucceededActionCreator(response))
	} catch (e) {
		yield put(failedActionCreator(e.message))
	}
}

interface SendMessagePayload {
	message: any;
	file:    any;
}

function* sendMessage({payload}: {payload: SendMessagePayload}) {
	try {
		const threadID = yield select(selectThreadID)
		const response = yield call(api.sendMessage,
				{message: payload.message, file: payload.file, threadID: threadID})

		if (response.error != "")
			yield put(failedActionCreator(response.error))
		else
			yield put(sendMessageSucceededActionCreator(response.message))
	} catch (e) {
		yield put(failedActionCreator(e.message))
	}
}

interface SetThreadMessagePayload {
	ID: number;
	[String]: any;
}

function* setThreadMessage({payload}: {payload: SetThreadMessagePayload}) {
	const threadID = yield select(selectThreadID)
	if (threadID !== payload.ID) {
		yield put(resetActionCreator())
	}
	yield put({
		type: __SET_THREAD_MESSAGE,
		payload,
	})
}

function* threadsSaga() {
	yield takeEvery(SET_THREAD_MESSAGE, setThreadMessage)
	yield takeLatest(FETCH_MESSAGES, fetchMessages)
	yield takeLatest(SEND_MESSAGE, sendMessage)
}

export {threadsSaga}