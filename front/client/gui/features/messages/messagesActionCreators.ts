import {
  UPDATE_MESSAGE,
  DELETE_MESSAGE,
  SEND_MESSAGE,
  FETCH_MESSAGES,
  SET_MESSAGE_FOR_EDIT,

  MESSAGES_FAILED,

  SEND_MESSAGE_SUCCEEDED,
  FETCH_MESSAGES_SUCCEEDED,
  UPDATE_MESSAGE_SUCCEEDED,
  DELETE_MESSAGE_SUCCEEDED,
} from './messagesActions'

export function messagesFailedActionCreator(payload) {
  return {
    type: MESSAGES_FAILED,
    payload,
  }
}

export function sendMessageActionCreator(message, file) {
  return {
    type: SEND_MESSAGE,
    payload: {message, file},
  }
}

export function sendMessageSucceededActionCreator(payload) {
  return {
    type: SEND_MESSAGE_SUCCEEDED,
    payload,
  }
}

export function fetchMessagesActionCreator(limit: number, offset: number, order: number) {
  return {
    type: FETCH_MESSAGES,
    payload: {limit, offset, order},
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

export function resetEditMessageActionCreator() {
  return setEditMessageActionCreator({})
}

export function deleteMessageActionCreator(payload) {
  return {
    type: DELETE_MESSAGE,
    payload,
  }
}

export function deleteMessageSucceededActionCreator(payload) {
  return {
    type: DELETE_MESSAGE_SUCCEEDED,
    payload,
  }
}