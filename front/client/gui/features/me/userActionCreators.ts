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

export const authActionCreator = (shouldSuccessRedirect, shouldFailRedirect) => ({
	type: AUTH,
	shouldSuccessRedirect,
	shouldFailRedirect,
})

export const authFailedActionCreator = (payload, shouldRedirect) => ({
	type: AUTH_FAILED,
	payload,
	shouldRedirect,
})

export const authSucceededActionCreator = (payload, shouldRedirect) => ({
	type: AUTH_SUCCEEDED,
	payload,
	shouldRedirect,
})

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