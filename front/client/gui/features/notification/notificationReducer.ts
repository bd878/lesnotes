import {SHOW_NOTIFICATION, HIDE_NOTIFICATION, SET_TIMER_NOTIFICATION} from './notificationActions';

const initialState = {
	timerID: null,
	isVisible: false,
	text: "",
}

export function notificationReducer(notificationState = initialState, action) {
	switch (action.type) {
	case SHOW_NOTIFICATION:
		return {
			...notificationState,
			isVisible: true,
			text: action.payload.text,
		}
	case HIDE_NOTIFICATION:
		return {
			...notificationState,
			isVisible: false,
			text: "",
			timerID: null,
		}
	case SET_TIMER_NOTIFICATION:
		return {
			...notificationState,
			timerID: action.payload.timerID,
		}
	default:
	}

	return notificationState
}
