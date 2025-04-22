import {SHOW_NOTIFICATION, HIDE_NOTIFICATION} from './notificationActions';

const initialState = {
	isVisible: false,
	text: "",
}

export function notificationReducer(notificationState = initialState, action) {
	switch (action.type) {
	case SHOW_NOTIFICATION:
	case HIDE_NOTIFICATION:
	default:
	}

	return notificationState
}
