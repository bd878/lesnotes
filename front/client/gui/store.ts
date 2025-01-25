import {createStore, combineReducers, applyMiddleware} from './third_party/redux'
import {messagesReducer, messagesSaga} from './features/messages'
import createSagaMiddleware from 'redux-saga'

const sagaMiddleware = createSagaMiddleware()

export default createStore(combineReducers({
  messages: messagesReducer,
}), {}, applyMiddleware(sagaMiddleware))

sagaMiddleware.run(messagesSaga)