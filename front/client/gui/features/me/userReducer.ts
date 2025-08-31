import {
	LOGOUT,
} from './userActions';
import models from '../../../api/models'

const initialState = {
	user: models.user(),
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
		}
	}
	}
	return userState
}