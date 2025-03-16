import {
	RESET,
	FAILED,
	SEND_MESSAGE,
	SEND_MESSAGE_SUCCEEDED,
	SET_THREAD_MESSAGE,
	FETCH_MESSAGES,
	FETCH_MESSAGES_SUCCEEDED,
} from './threadsActions';

const initialState = {
	list: [],
	threadID: 0,
	message: {},
	isLastPage: false,
	loading: false,
	error: "",
}

export function threadsReducer(state = initialState, action) {
	switch (action.type) {
		case FAILED: {
			return {
				...state,
				error: action.payload,
				loading: false,
			}
		}
		case FETCH_MESSAGES: {
			return {
				...state,
				errors: "",
				loading: true,
			}
		}
		case FETCH_MESSAGES_SUCCEEDED: {
			return {
				...state,
				list: [ ...action.payload.messages, ...state.list ],
				isLastPage: action.payload.isLastPage,
				loading: false,
				error: "",
			}
		}
		case SET_THREAD_MESSAGE: {
			return {
				...state,
				threadID: action.payload.ID,
				message: action.payload,
			}
		}
		case RESET: {
			return {
				...state,
				threadID: 0,
				message: {},
				list: [],
			}
		}
		case SEND_MESSAGE: {
			return {
				...state,
				loading: true,
				error: "",
			}
		}
		case SEND_MESSAGE_SUCCEEDED: {
			return {
				...state,
				list: [ ...state.list, action.payload ],
			}
		}
	}
	/* init */
	return state
}