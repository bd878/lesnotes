import {takeLatest,put,call,select} from 'redux-saga/effects'
import {FETCH_MESSAGES} from './threadsActions'
import {selectThreadID} from './threadsSelectors'
import {
	messagesFailedActionCreator,
	fetchMessagesSucceededActionCreator,
	messagesFailedActionCreator,
} from './threadsActionCreators'
import api from '../../api'

interface FetchMessagesPayload {
	limit:  number;
	offset: number;
	order:  number;
}

function* fetchMessages({payload}: {payload: FetchMessagesPayload}) {
	try {
		// TODO: check of threadID is not set, fetch is valid with threadID
		const threadID = yield select(selectThreadID)
		const response = yield call(api.loadMessages,
			{limit: payload.limit, offset: payload.offset,
			order: payload.order, threadID: threadID})

		response.messages.reverse();
		if (response.error != "")
			yield put(messagesFailedActionCreator(response.error))
		else
			yield put(fetchMessagesSucceededActionCreator(response))
	} catch (e) {
		yield put(messagesFailedActionCreator(e.message))
	}
}

function* threadsSaga() {
	yield takeLatest(FETCH_MESSAGES, fetchMessages)
}

export {threadsSaga}