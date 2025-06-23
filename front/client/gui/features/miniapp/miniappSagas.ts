import {takeLatest, put, call, fork, spawn, select} from 'redux-saga/effects'
import {VALIDATE_INIT_DATA} from './miniappActions'
import {
	validateInitDataSucceededActionCreator,
	miniappFailedActionCreator,
} from './miniappActionCreators';
import * as is from '../../../third_party/is'
import api from '../../../api'

function* validateInitData({payload}) {
	try {
		let result = yield call(api.validateMiniappData, payload)

		if (result.ok) {
			yield put(validateInitDataSucceededActionCreator(result.token))
		} else {
			yield put(miniappFailedActionCreator(result.error + ": " + result.explain))
			console.log(JSON.stringify(result))
		}
	} catch (e) {
		yield put(miniappFailedActionCreator(e.message))
		console.error(e)
	}
}

function* miniappSaga() {
	yield takeLatest(VALIDATE_INIT_DATA, validateInitData)
}

export {miniappSaga}