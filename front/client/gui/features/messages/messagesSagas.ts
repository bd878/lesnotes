import {takeLatest,put,call,select} from 'redux-saga/effects'
import {UPDATE_MESSAGE, FETCH_MESSAGES, SEND_MESSAGE} from './messagesActions'
import {
  updateMessageActionCreator,
  updateMessageSucceededActionCreator,
  updateMessageFailedActionCreator,
  sendMessageFailedActionCreator,
  fetchMessagesFailedActionCreator,
  fetchMessagesSucceededActionCreator,
  pushBackMessagesActionCreator,
  appendMessagesActionCreator,
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

    if (response.error != "")
      yield put(sendMessageFailedActionCreator(response.error))
    else
      yield put(appendMessagesActionCreator([response.message]))
  } catch (e) {
    yield put(sendMessageFailedActionCreator(e.message))
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
      yield put(updateMessageFailedActionCreator(response.error))
    else
      yield put(updateMessageSucceededActionCreator(messages))
  } catch (e) {
    yield put(updateMessageFailedActionCreator(e.message))
  }
}

function* messagesSaga() {
  yield takeLatest(UPDATE_MESSAGE, updateMessage)
  yield takeLatest(FETCH_MESSAGES, fetchMessages)
  yield takeLatest(SEND_MESSAGE, sendMessage)
}

export {messagesSaga}