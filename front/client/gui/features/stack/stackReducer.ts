import {
	OPEN_THREAD,
	CLOSE_THREAD,
	DESTROY_THREAD,

	UPDATE_MESSAGE,
	DELETE_MESSAGE,
	SEND_MESSAGE,
	COPY_MESSAGE,
	FETCH_MESSAGES,
	SET_MESSAGE_FOR_EDIT,

	MESSAGES_FAILED,

	SEND_MESSAGE_SUCCEEDED,
	FETCH_MESSAGES_SUCCEEDED,
	UPDATE_MESSAGE_SUCCEEDED,
	DELETE_MESSAGE_SUCCEEDED,
} from './stackActions';
import * as is from '../../../third_party/is'

const thread = {
	ID: 0,
	list: [],
	messageForEdit: {},
	isLastPage: false,
	loading: false,
	error: "",
}

const initialState = {
	stack: [structuredClone(thread)],
}

function messageReducer(messagesState = thread, action) {
	switch (action.type) {
	case MESSAGES_FAILED: {
		return {
			...messagesState,
			error: action.payload,
			loading: false,
			messageForEdit: {},
		}
	}
	case FETCH_MESSAGES: {
		return {
			...messagesState,
			errors: "",
			loading: true,
		}
	}
	case FETCH_MESSAGES_SUCCEEDED: {
		return {
			...messagesState,
			list: [ ...action.payload.messages, ...messagesState.list ],
			isLastPage: action.payload.isLastPage,
			loading: false,
			error: "",
		}
	}
	case SEND_MESSAGE: {
		return {
			...messagesState,
			loading: true,
			error: "",
		}
	}
	case COPY_MESSAGE: {
		return messagesState
	}
	case SEND_MESSAGE_SUCCEEDED: {
		return {
			...messagesState,
			list: [ ...messagesState.list, action.payload ],
		}
	}
	case UPDATE_MESSAGE: {
		return {
			...messagesState,
			loading: true,
			error: "",
		}
	}
	case SET_MESSAGE_FOR_EDIT: {
		return {
			...messagesState,
			messageForEdit: action.payload,
		}
	}
	case DELETE_MESSAGE: {
		return {
			...messagesState,
			loading: true,
			error: "",
		}
	}
	case DELETE_MESSAGE_SUCCEEDED: {
		return {
			...messagesState,
			list: [ ...action.payload ],
			loading: false,
			error: "",
		}
	}
	case UPDATE_MESSAGE: {
		return {
			...messagesState,
			loading: true,
			error: "",
		}
	}
	case UPDATE_MESSAGE_SUCCEEDED: {
		return {
			...messagesState,
			list: [ ...action.payload ],
			loading: false,
			error: "",
			messageForEdit: {},
		}
	}
	case SET_MESSAGE_FOR_EDIT: {
		return {
			...messagesState,
			messageForEdit: action.payload,
		}
	}
	}
	return messagesState
}

export function stackReducer(stackState = initialState, action) {
	switch (action.type) {
		case OPEN_THREAD: {
			const nextStack = stackState.stack.slice(0, action.payload.index+1)

			const nextThread = structuredClone(thread)
			nextThread.ID = action.payload.threadID
			nextStack.push(nextThread)

			return {
				...stackState,
				stack: nextStack,
			}
		}
		case CLOSE_THREAD: {
			const index = action.payload.index
			return {
				...stackState,
				stack: stackState.stack.slice(0, index+1),
			}
		}
		case DESTROY_THREAD: {
			return {
				...stackState,
				stack: stackState.stack.slice(0, action.payload.index),
			}
		}
	}

	if (is.notUndef(action.index)) {
		let messageState = stackState.stack[action.index]
		if (is.notUndef(messageState)) {
			stackState.stack[action.index] = messageReducer(messageState, action)
		} else {
			console.error("cannot find Thread by index", action.index, action.type)
		}
		return {
			...stackState
		}
	}

	/*init*/
	return stackState
}
