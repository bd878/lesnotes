import {
  FETCH_MESSAGES,
  FETCH_MESSAGES_FAILED,
  FETCH_MESSAGES_SUCCEEDED,
  APPEND_MESSAGES,
  UPDATE_MESSAGE,
  UPDATE_MESSAGE_FAILED,
  UPDATE_MESSAGE_SUCCEEDED,
  PUSH_BACK_MESSAGES,
  SET_MESSAGE_FOR_EDIT,
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
        isLastPage: action.payload.isLastPage,
        loading: false,
        error: "",
      }
    }
    case FETCH_MESSAGES_FAILED: {
      return {
        ...messagesState,
        loading: false,
        errors: action.payload,
      }
    }
    case APPEND_MESSAGES: {
      return {
        ...messagesState,
        list: [ ...messagesState.list, ...action.payload ],
      }
    }
    case PUSH_BACK_MESSAGES: {
      return {
        ...messagesState,
        list: [ ...action.payload, ...messagesState.list ],
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
    case UPDATE_MESSAGE_FAILED: {
      return {
        ...messagesState,
        loading: false,
        error: action.payload,
      }
    }
    case SET_MESSAGE_FOR_EDIT: {
      return {
        ...messagesState,
        messageForEdit: action.payload,
      }
    }
  }
  /* init */
  return messagesState
}