import {createStore, combineReducers, applyMiddleware} from '../third_party/redux'
import {all} from 'redux-saga/effects'
import {userReducer, userSaga} from './features/me'
import {stackReducer, stackSaga} from './features/stack'
import createSagaMiddleware from 'redux-saga'

const sagaMiddleware = createSagaMiddleware()

export default ({
	browser = "",
	isMobile = false,
	isDesktop = true,
} = {}) => {
	const store = createStore(combineReducers({
		me: userReducer,
		stack: stackReducer,
	}), {
		me: {
			browser,
			isMobile,
			isDesktop,
		},
	}, applyMiddleware(sagaMiddleware))

	sagaMiddleware.run(function* rootSaga() {
		yield all([
			userSaga(),
			stackSaga(),
		])
	})

	return store
}