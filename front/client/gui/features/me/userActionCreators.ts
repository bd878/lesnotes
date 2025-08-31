import {
	LOGOUT,
	AUTH,
	AUTH_FAILED,
	AUTH_SUCCEEDED,
	WILL_REDIRECT,
	RESET_REDIRECT,
} from './userActions'

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