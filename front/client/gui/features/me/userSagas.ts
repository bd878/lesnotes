import i18n from '../../../i18n';
import {takeLatest,put,call} from 'redux-saga/effects'
import {
	LOGIN,
	LOGOUT,
	AUTH,
	REGISTER,
	AUTH_FAILED,
	AUTH_SUCCEEDED,
	LOGIN_SUCCEEDED,
	REGISTER_SUCCEEDED,
} from './userActions'
import {
	logoutActionCreator,
	authSucceededActionCreator,
	authFailedActionCreator,
	loginSucceededActionCreator,
	loginFailedActionCreator,
	registerFailedActionCreator,
	registerSucceededActionCreator,
	willRedirectActionCreator,
	resetRedirectActionCreator,
} from './userActionCreators';
import * as is from '../../../third_party/is'
import api from '../../../api'

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

function* redirectHome({shouldRedirect = true}) {
	if (location && !location.pathname.includes("/home") && shouldRedirect)
		setTimeout(() => {location.href = "/home"})
	else
		yield put(resetRedirectActionCreator())
}

interface RegisterPayload {
	name:   string;
	password: string;
}

function* register({payload}: {payload: RegisterPayload}) {
	try {
		const response = yield call(api.register,
			payload.name, payload.password)

		if (is.notEmpty(response.error))
			yield put(registerFailedActionCreator(response.error))
		else
			yield put(registerSucceededActionCreator(response))
	} catch (e) {
		yield put(registerFailedActionCreator(e.message))
	}
}

function* auth({shouldSuccessRedirect = true, shouldFailRedirect = true}) {
	try {
		const response = yield call(api.auth)

		yield put(willRedirectActionCreator())
		if (is.notEmpty(response.error))
			yield put(authFailedActionCreator(response.error, shouldFailRedirect))
		else
			yield put(authSucceededActionCreator(response, shouldSuccessRedirect))
	} catch (e) {
		yield put(authFailedActionCreator(e.message, shouldFailRedirect))
	}
}

function* redirectLogin({shouldRedirect = true}) {
	if (location && !location.pathname.includes("/login") && shouldRedirect)
		setTimeout(() => {location.href = "/login"})
	else
		yield put(resetRedirectActionCreator())
}

function* logout() {
	try {
		yield call(api.logout)
		setTimeout(() => {location.href = "/login"})
	} catch (e) {
		console.error(i18n("error_occured"), e);
	}
}

function* userSaga() {
	yield takeLatest(LOGOUT, logout)
	yield takeLatest(AUTH, auth)
	yield takeLatest(LOGIN, login)
	yield takeLatest(REGISTER, register)
	yield takeLatest(AUTH_FAILED, redirectLogin)
	yield takeLatest(AUTH_SUCCEEDED, redirectHome)
	yield takeLatest(LOGIN_SUCCEEDED, redirectHome)
	yield takeLatest(REGISTER_SUCCEEDED, redirectHome)
}

export {userSaga}