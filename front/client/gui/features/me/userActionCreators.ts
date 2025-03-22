import {
	LOGOUT,
	LOGIN,
	LOGIN_SUCCEEDED,
	LOGIN_FAILED,
	REGISTER,
	REGISTER_FAILED,
	REGISTER_SUCCEEDED,
	AUTH,
	AUTH_FAILED,
	AUTH_SUCCEEDED,
	WILL_REDIRECT,
	RESET_REDIRECT,
} from './userActions'

export function registerActionCreator(name, password) {
	return {
		type: REGISTER,
		payload: {name, password},
	}
}

export function registerSucceededActionCreator(payload) {
	return {
		type: REGISTER_SUCCEEDED,
		payload,
	}
}

export function registerFailedActionCreator(payload) {
	return {
		type: REGISTER_FAILED,
		payload,
	}
}

export function loginActionCreator(name, password) {
	return {
		type: LOGIN,
		payload: {name, password},
	}
}

export function loginSucceededActionCreator(payload) {
	return {
		type: LOGIN_SUCCEEDED,
		payload,
	}
}

export function loginFailedActionCreator(payload) {
	return {
		type: LOGIN_FAILED,
		payload,
	}
}

export function authActionCreator(payload = {}) {
	return {
		type: AUTH,
		payload: payload,
	}
}

export function authFailedActionCreator(payload) {
	return {
		type: AUTH_FAILED,
		payload,
	}
}

export function authSucceededActionCreator(payload) {
	return {
		type: AUTH_SUCCEEDED,
		payload,
	}
}

export function logoutActionCreator(payload = {}) {
	return {
		type: LOGOUT,
		payload,
	}
}

export function willRedirectActionCreator(payload = {}) {
	return {
		type: WILL_REDIRECT,
		payload,
	}
}

export function resetRedirectActionCreator(payload = {}) {
	return {
		type: RESET_REDIRECT,
		payload,
	}
}