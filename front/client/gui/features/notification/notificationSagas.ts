import {takeLatest, put, call, fork, spawn, select} from 'redux-saga/effects'
import {SHOW_NOTIFICATION, HIDE_NOTIFICATION} from './notificationActions'
import {
	hideNotificationActionCreator,
	setNotificationTimerActionCreator,
} from './notificationActionCreators';
import * as is from '../../../third_party/is'
import {selectTimerID} from './notificationSelectors';

const NOTIFICATION_TIMEOUT_SEC = 3000

function* showNotification({payload}) {
	let timerID = yield select(selectTimerID)
	if (is.notEmpty(timerID))
		clearTimeout(timerID)

	yield spawn(function* dispatchHideNotification() {
		let timeout = {fn: () => {}}
		let timerID = setTimeout(() => {timeout.fn()}, NOTIFICATION_TIMEOUT_SEC)
		yield put(setNotificationTimerActionCreator({timerID}))
		yield call(() => new Promise(resolve => timeout.fn = resolve))
		yield put(hideNotificationActionCreator())
	})
}

function* hideNotification({payload}) {}

function* notificationSaga() {
	yield takeLatest(SHOW_NOTIFICATION, showNotification)
	yield takeLatest(HIDE_NOTIFICATION, hideNotification)
}

export {notificationSaga}