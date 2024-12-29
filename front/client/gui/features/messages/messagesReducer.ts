import {
  APPEND_MESSAGES,
  PUSH_BACK_MESSAGES,
} from './messagesActions';

const initialState = {
  list: [],
}

export function messagesReducer(messagesState = initialState, action) {
  switch (action.type) {
    case APPEND_MESSAGES: {
      return {
        list: [ ...messagesState.list, ...action.payload ],
      }
    }
    case PUSH_BACK_MESSAGES: {
      return {
        list: [ ...action.payload, ...messagesState.list ],
      }
    }
  }
  /* init */
  return messagesState
}