import {createStore, combineReducers} from './third_party/redux'
import {messagesReducer} from './features/messages'

export const store = createStore(combineReducers({
  messages: messagesReducer,
}), {})