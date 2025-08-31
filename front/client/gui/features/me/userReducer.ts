import {
	LOGOUT,
	AUTH,
	AUTH_FAILED,
	AUTH_SUCCEEDED,
	WILL_REDIRECT,
	RESET_REDIRECT,
} from './userActions';
import models from '../../../api/models'

const initialState = {
	user: models.user(),
	isAuth: false,
	loading: false,
	error: "",
	willRedirect: true,
	browser: "",
	isMobile: false,
	isDesktop: true,
	isMiniapp: false,
}

export function userReducer(userState = initialState, action) {
	switch (action.type) {
	case LOGOUT: {
		return {
			...userState,
			user: models.user(),
			isAuth: false,
			loading: false,
			error: "",
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
	case WILL_REDIRECT: {
		return {
			...userState,
			willRedirect: true,
		}
	}
	case RESET_REDIRECT: {
		return {
			...userState,
			willRedirect: false,
		}
	}
	}
	return userState
}