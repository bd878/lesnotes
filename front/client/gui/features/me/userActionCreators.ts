import {
	LOGOUT,
} from './userActions'

export function logoutActionCreator(payload = {}) {
	return {
		type: LOGOUT,
		payload,
	}
}
