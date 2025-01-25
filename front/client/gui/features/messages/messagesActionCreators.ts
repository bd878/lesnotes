import {
  APPEND_MESSAGES,
  PUSH_BACK_MESSAGES,
  FETCH_MESSAGES_SUCCEEDED,
  FETCH_MESSAGES,
  FETCH_MESSAGES_FAILED,
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

interface FetchMessagesPayload {
  limit:  number;
  offset: number;
  order:  number;
}

export function fetchMessagesActionCreator(payload: FetchMessagesPayload) {
  return {
    type: FETCH_MESSAGES,
    payload,
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