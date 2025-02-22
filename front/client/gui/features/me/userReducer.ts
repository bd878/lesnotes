import {
  LOGIN,
  LOGIN_FAILED,
  LOGIN_SUCCEEDED,
  REGISTER,
  REGISTER_FAILED,
  REGISTER_SUCCEEDED,
  AUTH,
  AUTH_FAILED,
  AUTH_SUCCEEDED,
} from './userActions';
import models from '../../api/models'

const initialState = {
  user: models.user(),
  isAuth: false,
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
  case AUTH: {
    return {
      ...userState,
      loading: true,
      error: "",
    }
  }
  case AUTH_FAILED: {
    return {
      ...userState,
      isAuth: false,
      loading: false,
      error: action.payload,
    }
  }
  case AUTH_SUCCEEDED: {
    return {
      ...userState,
      isAuth: true,
      loading: false,
      error: "",
      user: action.payload.user,
    }
  }
  }
  return userState
}