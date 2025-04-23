import {SHOW_NOTIFICATION, HIDE_NOTIFICATION, SET_TIMER_NOTIFICATION} from './notificationActions'

export const showNotificationActionCreator = payload => ({
	type: SHOW_NOTIFICATION,
	payload,
})

export const hideNotificationActionCreator = payload => ({
	type: HIDE_NOTIFICATION,
	payload,
})

export const setNotificationTimerActionCreator = payload => ({
	type: SET_TIMER_NOTIFICATION,
	payload,
})