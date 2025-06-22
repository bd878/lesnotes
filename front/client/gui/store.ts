import {createStore as createReduxStore, combineReducers, applyMiddleware} from '../third_party/redux'
import {all} from 'redux-saga/effects'
import {userReducer, userSaga} from './features/me'
import {stackReducer, stackSaga} from './features/stack'
import {notificationReducer, notificationSaga} from './features/notification'
import {miniappReducer, miniappSaga} from './features/miniapp'
import createSagaMiddleware from 'redux-saga'

let instance = null

const sagaMiddleware = createSagaMiddleware()

export default function createStore({
	browser = "",
	isMobile = false,
	isDesktop = true,
} = {}) {
	if (instance == null) {
		instance = createReduxStore(combineReducers({
			me: userReducer,
			stack: stackReducer,
			notification: notificationReducer,
			miniapp: miniappReducer,
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
				notificationSaga(),
				miniappSaga(),
			])
		})
	}

	return instance
}
