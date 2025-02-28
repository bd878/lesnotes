import {
  DELETE_MESSAGE,
  FETCH_MESSAGES,
  UPDATE_MESSAGE,
  SEND_MESSAGE,
  SET_MESSAGE_FOR_EDIT,

  MESSAGES_FAILED,

  FETCH_MESSAGES_SUCCEEDED,
  SEND_MESSAGE_SUCCEEDED,
  UPDATE_MESSAGE_SUCCEEDED,
  DELETE_MESSAGE_SUCCEEDED,
} from './messagesActions';

const initialState = {
  list: [],
  messageForEdit: {},
  isLastPage: false,
  loading: false,
  error: "",
}

export function messagesReducer(messagesState = initialState, action) {
  switch (action.type) {
    case MESSAGES_FAILED: {
      return {
        ...messagesState,
        error: action.payload,
        loading: false,
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
    case UPDATE_MESSAGE_SUCCEEDED: {
      return {
        ...messagesState,
        list: [ ...action.payload ],
        loading: false,
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
  }
  /* init */
  return messagesState
}