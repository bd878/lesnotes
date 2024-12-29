import {
  APPEND_MESSAGES,
  PUSH_BACK_MESSAGES
} from './messagesActions'

export function addMessagesActionCreator(payload) {
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