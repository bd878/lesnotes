import {VALIDATE_INIT_DATA, VALIDATE_INIT_DATA_SUCCEEDED, MINIAPP_FAILED} from './miniappActions';

const initialState = {
	loading: false,
	valid: false,
	token: "",
	error: "",
}

export function miniappReducer(miniappState = initialState, action) {
	switch (action.type) {
	case MINIAPP_FAILED: {
		return {
			...miniappState,
			error: action.payload,
			loading: false,
		}
	}
	case VALIDATE_INIT_DATA: {
		return {
			...miniappState,
			loading: true,
			valid: false,
		}
	}
	case VALIDATE_INIT_DATA_SUCCEEDED: {
		return {
			...miniappState,
			loading: false,
			valid: true,
			token: action.payload,
			error: "",
		}
	}
	default:
	}

	return miniappState
}
