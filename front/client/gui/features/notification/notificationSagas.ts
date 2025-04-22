import {takeLatest} from 'redux-saga/effects'
import {SHOW_NOTIFICATION, HIDE_NOTIFICATION} from './notificationActions'

function* showNotification({payload}) {}

function* hideNotification({payload}) {}

function* notificationSaga() {
	yield takeLatest(SHOW_NOTIFICATION, showNotification)
	yield takeLatest(HIDE_NOTIFICATION, hideNotification)
}

export {notificationSaga}