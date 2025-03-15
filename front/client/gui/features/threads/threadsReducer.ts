import {
	SET_THREAD_MESSAGE,
	FETCH_MESSAGES,
	MESSAGES_FAILED,
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
		case MESSAGES_FAILED: {
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
	}
	/* init */
	return state
}