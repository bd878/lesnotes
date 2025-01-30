import {
  SEND_MESSAGE,
  SEND_MESSAGE_FAILED,
  SEND_MESSAGE_SUCCEEDED,
  APPEND_MESSAGES,
  PUSH_BACK_MESSAGES,
  FETCH_MESSAGES_SUCCEEDED,
  FETCH_MESSAGES,
  FETCH_MESSAGES_FAILED,
} from './messagesActions'

export function sendMessageActionCreator(message, file) {
  console.log("action creator", "message=", message, "file=", file)
  return {
    type: SEND_MESSAGE,
    payload: {message, file},
  }
}

export function sendMessageFailedActionCreator(payload) {
  return {
    type: SEND_MESSAGE_FAILED,
    payload,
  }
}

export function sendMessageSucceededActionCreator(payload) {
  return {
    type: SEND_MESSAGE_SUCCEEDED,
    payload,
  }
}

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

export function fetchMessagesActionCreator(limit: number, offset: number, order: number) {
  return {
    type: FETCH_MESSAGES,
    payload: {limit, offset, order},
  }
}

export function fetchMessagesFailedActionCreator(payload) {
  return {
    type: FETCH_MESSAGES_FAILED,
    payload,
  }
}

export function fetchMessagesSucceededActionCreator(payload) {
  return {
    type: FETCH_MESSAGES_SUCCEEDED,
    payload,
  }
}