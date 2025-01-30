import {
  LOGIN,
  LOGIN_FAILED,
  LOGIN_SUCCEEDED,
  REGISTER,
  REGISTER_FAILED,
  REGISTER_SUCCEEDED,
} from './userActions';

const initialState = {
  user: {
    name: "",
    id: "",
  },
  loading: false,
  error: "",
}

export function userReducer(userState = initialState, action) {
  switch (action.type) {
    case LOGIN: {
      return {
        ...userState,
        loading: true,
        error: ""
      }
    }
    case REGISTER: {
      return {
        ...userState,
        loading: true,
        error: ""
      }
    }
    case LOGIN_FAILED: {
      return {
        ...userState,
        loading: false,
        error: action.payload,
      }
    }
    case LOGIN_SUCCEEDED: {
      return {
        ...userState,
        loading: false,
        error: "",
      }
    }
    case REGISTER_FAILED: {
      return {
        ...userState,
        loading: false,
        error: action.payload,
      }
    }
    case REGISTER_SUCCEEDED: {
      return {
        ...userState,
        loading: false,
        error: ""
      }
    }
  }
  return userState
}