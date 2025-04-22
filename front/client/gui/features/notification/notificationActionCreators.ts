import {SHOW_NOTIFICATION, HIDE_NOTIFICATION} from './notificationActions'

export const showNotificationActionCreator = payload => ({
	type: SHOW_NOTIFICATION,
	payload,
})

export const hideNotificationActionCreator = payload => ({
	type: HIDE_NOTIFICATION,
	payload,
})