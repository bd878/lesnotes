import {takeLatest,put,call} from 'redux-saga/effects'
import {
  LOGIN,
  AUTH,
  REGISTER,
  AUTH_SUCCEEDED,
  LOGIN_SUCCEEDED,
  REGISTER_SUCCEEDED,
} from './userActions'
import {
  authSucceededActionCreator,
  authFailedActionCreator,
  loginSucceededActionCreator,
  loginFailedActionCreator,
  registerFailedActionCreator,
  registerSucceededActionCreator,
} from './userActionCreators';
import api from '../../api'

interface LoginPayload {
  name:   string;
  password: string;
}

function* login({payload}: {payload: LoginPayload}) {
  try {
    const response = yield call(api.login,
      payload.name, payload.password)

    if (response.error !== "")
      yield put(loginFailedActionCreator(response.error))
    else
      yield put(loginSucceededActionCreator(response))
  } catch (e) {
    yield put(loginFailedActionCreator(e.message))
  }
}

function* redirectHome() {
  setTimeout(() => {location.href = "/home"}, 0)
}

interface RegisterPayload {
  name:   string;
  password: string;
}

function* register({payload}: {payload: RegisterPayload}) {
  try {
    const response = yield call(api.register,
      payload.name, payload.password)

    if (response.error !== "")
      yield put(registerFailedActionCreator(response.error))
    else
      yield put(registerSucceededActionCreator(response))
  } catch (e) {
    yield put(registerFailedActionCreator(e.message))
  }
}

interface AuthPayload {}

function* auth({payload}) {
  try {
    const response = yield call(api.auth)

    if (response.error !== "")
      yield put(authFailedActionCreator(response.error))
    else
      yield put(authSucceededActionCreator(response))
  } catch (e) {
    yield put(authFailedActionCreator(e.message))
  }
}

function* userSaga() {
  yield takeLatest(AUTH, auth)
  yield takeLatest(LOGIN, login)
  yield takeLatest(REGISTER, register)
  yield takeLatest(LOGIN_SUCCEEDED, redirectHome)
  yield takeLatest(REGISTER_SUCCEEDED, redirectHome)
}

export {userSaga}