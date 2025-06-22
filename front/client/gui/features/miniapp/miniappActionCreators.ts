import {VALIDATE_INIT_DATA, VALIDATE_INIT_DATA_SUCCEEDED, MINIAPP_FAILED} from './miniappActions'

export const validateInitDataActionCreator = payload => ({
	type: VALIDATE_INIT_DATA,
	payload,
})

export const validateInitDataSucceededActionCreator = payload => ({
	type: VALIDATE_INIT_DATA_SUCCEEDED,
	payload,
})

export const miniappFailedActionCreator = payload => ({
	type: MINIAPP_FAILED,
	payload,
})