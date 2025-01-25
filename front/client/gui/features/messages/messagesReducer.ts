import {
  FETCH_MESSAGES,
  FETCH_MESSAGES_FAILED,
  FETCH_MESSAGES_SUCCEEDED,
  APPEND_MESSAGES,
  PUSH_BACK_MESSAGES,
} from './messagesActions';

const initialState = {
  list: [],
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
  }
  /* init */
  return messagesState
}