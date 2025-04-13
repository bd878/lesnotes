import {createStore, combineReducers, applyMiddleware} from '../third_party/redux'
import {all} from 'redux-saga/effects'
import {messagesReducer, messagesSaga} from './features/messages'
import {threadsReducer, threadsSaga} from './features/threads'
import {userReducer, userSaga} from './features/me'
import createSagaMiddleware from 'redux-saga'

const sagaMiddleware = createSagaMiddleware()

export default ({
	browser = "",
	isMobile = false,
	isDesktop = true,
} = {}) => {
	const store = createStore(combineReducers({
		messages: messagesReducer,
		threads: threadsReducer,
		me: userReducer,
	}), {
		me: {
			browser,
			isMobile,
			isDesktop,
		},
	}, applyMiddleware(sagaMiddleware))

	sagaMiddleware.run(function* rootSaga() {
		yield all([
			threadsSaga(),
			userSaga(),
			messagesSaga(),
		])
	})

	return store
}
