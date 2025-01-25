import {takeLatest,put,call} from 'redux-saga/effects'
import {FETCH_MESSAGES, SEND_MESSAGE} from './messagesActions'
import {
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

function* fetchMessages(action: FetchMessagesPayload) {
  try {
    const response = yield call(api.loadMessages,
      action.limit, action.offset, action.order)

    if (response.error != "") {
      yield put(fetchMessagesFailedActionCreator(response.error))
    } else {
      yield put(fetchMessagesSucceededActionCreator(response))
      yield put(pushBackMessagesActionCreator(response.messages))
    }
  } catch (e) {
    yield put(fetchMessagesFailedActionCreator(e.message))
  }
}

function* messagesSaga() {
  yield takeLatest(FETCH_MESSAGES, fetchMessages)
}

export {messagesSaga}