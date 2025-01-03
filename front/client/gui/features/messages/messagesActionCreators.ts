import {
  APPEND_MESSAGES,
  PUSH_BACK_MESSAGES
} from './messagesActions'

export function appendMessagesActionCreator(payload) {
  return {
    type: APPEND_MESSAGES,
    payload,
  }
}

export function pushBackMessagesActionCreator(payload) {
  return {
    type: PUSH_BACK_MESSAGES,
    payload,
  }
}