import {createStore, combineReducers, applyMiddleware} from './third_party/redux'
import {all} from 'redux-saga/effects'
import {messagesReducer, messagesSaga} from './features/messages'
import {userReducer, userSaga} from './features/me'
import createSagaMiddleware from 'redux-saga'

const sagaMiddleware = createSagaMiddleware()

export default createStore(combineReducers({
  messages: messagesReducer,
  me: userReducer,
}), {}, applyMiddleware(sagaMiddleware))

sagaMiddleware.run(function* rootSaga() {
  yield all([
    userSaga(),
    messagesSaga(),
  ])
})