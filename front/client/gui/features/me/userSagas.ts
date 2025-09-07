import i18n from '../../../i18n';
import {takeLatest,put,call} from 'redux-saga/effects'
import {
	LOGOUT,
} from './userActions'
import {
	logoutActionCreator,
} from './userActionCreators';
import * as is from '../../../third_party/is'
import api from '../../../api'

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
}

export {userSaga}