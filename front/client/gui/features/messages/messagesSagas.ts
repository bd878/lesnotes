import {takeLatest,put,call} from 'redux-saga/effects'
import {FETCH_MESSAGES, SEND_MESSAGE} from './messagesActions'
import {
  sendMessageFailedActionCreator,
  fetchMessagesFailedActionCreator,
  fetchMessagesSucceededActionCreator,
  pushBackMessagesActionCreator,
} from './messagesActionCreators'
import api from '../../api'

interface FetchMessagesPayload {
  limit:  number;
  offset: number;
  order:  number;
}

function* fetchMessages({payload}: {payload: FetchMessagesPayload}) {
  try {
    const response = yield call(api.loadMessages,
      payload.limit, payload.offset, payload.order)

    if (response.error != "") {
      yield put(fetchMessagesFailedActionCreator(response.error))
    } else {
      response.messages.reverse();
      yield put(fetchMessagesSucceededActionCreator(response))
      yield put(pushBackMessagesActionCreator(response.messages))
    }
  } catch (e) {
    yield put(fetchMessagesFailedActionCreator(e.message))
  }
}

interface SendMessagePayload {
  message: any;
  file:    any;
}

function* sendMessage({payload}: {payload: SendMessagePayload}) {
  try {
    const response = yield call(api.sendMessage,
        payload.message, payload.file)

    if (response.error != "") {
      yield put(sendMessageFailedActionCreator(response.error))
    } else {
      yield put(pushBackMessagesActionCreator([response.message]))
    }
  } catch (e) {
    yield put(sendMessageFailedActionCreator(e.message))
  }
}

function* messagesSaga() {
  yield takeLatest(FETCH_MESSAGES, fetchMessages)
  yield takeLatest(SEND_MESSAGE, sendMessage)
}

export {messagesSaga}