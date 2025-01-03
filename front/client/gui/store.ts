import {createStore, combineReducers} from './third_party/redux'
import {messagesReducer} from './features/messages'

export default createStore(combineReducers({
  messages: messagesReducer,
}), {})