import {
  UPDATE_MESSAGE,
  UPDATE_MESSAGE_FAILED,
  UPDATE_MESSAGE_SUCCEEDED,
  SEND_MESSAGE,
  SEND_MESSAGE_FAILED,
  SEND_MESSAGE_SUCCEEDED,
  APPEND_MESSAGES,
  PUSH_BACK_MESSAGES,
  FETCH_MESSAGES_SUCCEEDED,
  FETCH_MESSAGES,
  FETCH_MESSAGES_FAILED,
  SET_MESSAGE_FOR_EDIT,
} from './messagesActions'

export function sendMessageActionCreator(message, file) {
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

export function updateMessageActionCreator(ID, text) {
  return {
    type: UPDATE_MESSAGE,
    payload: {ID, text},
  }
}

export function updateMessageFailedActionCreator(payload) {
  return {
    type: UPDATE_MESSAGE_FAILED,
    payload,
  }
}

export function updateMessageSucceededActionCreator(payload) {
  return {
    type: UPDATE_MESSAGE_SUCCEEDED,
    payload,
  }
}

export function setEditMessageActionCreator(payload) {
  return {
    type: SET_MESSAGE_FOR_EDIT,
    payload,
  }
}