import {takeLatest,put,call,select} from 'redux-saga/effects'
import {
	UPDATE_MESSAGE,
	FETCH_MESSAGES,
	SEND_MESSAGE,
	DELETE_MESSAGE,
	COPY_MESSAGE,
} from './stackActions'
import {
	messagesFailedActionCreator,
	fetchMessagesSucceededActionCreator,
	sendMessageSucceededActionCreator,
	updateMessageSucceededActionCreator,
	deleteMessageSucceededActionCreator,
} from './stackActionCreators'
import * as is from '../../../third_party/is'
import {selectMessages} from './stackSelectors';
import {selectBrowser, selectIsMobile, selectIsDesktop} from '../me'
import api from '../../../api'

interface FetchMessagesPayload {
	limit:  number;
	offset: number;
	order:  number;
}

function* fetchMessages({index, payload}: {payload: FetchMessagesPayload}) {
	try {
		const response = yield call(api.loadMessages,
			{limit: payload.limit, offset: payload.offset, order: payload.order})

		response.messages.reverse();
		if (response.error != "")
			yield put(messagesFailedActionCreator(index)(response.error))
		else
			yield put(fetchMessagesSucceededActionCreator(index)(response))
	} catch (e) {
		yield put(messagesFailedActionCreator(index)(e.message))
	}
}

function* sendMessage({index, payload}) {
	try {
		let response
		if (is.notUndef(payload.file)) {
			response = yield call(api.uploadFile, payload.file)
			if (response.error != "") {
				yield put(messagesFailedActionCreator(index)(response.error))
				return
			}

			response = yield call(api.sendMessage, {text: payload.text, fileID: response.ID})
		} else {
			response = yield call(api.sendMessage, {text: payload.text})
		}

		if (response.error != "")
			yield put(messagesFailedActionCreator(index)(response.error))
		else
			yield put(sendMessageSucceededActionCreator(index)(response.message))
	} catch (e) {
		yield put(messagesFailedActionCreator(index)(e.message))
	}
}

function* updateMessage({index, payload}) {
	try {
		const response = yield call(api.updateMessage, payload.ID, payload.text)

		const messages = yield select(selectMessages(index))
		let idx = messages.findIndex(({ID}) => ID === payload.ID)
		if (idx !== -1)
			messages[idx] = {
				...messages[idx],
				ID: response.ID,
				text: payload.text,
				updateUTCNano: response.updateUTCNano,
			}

		if (response.error !== "")
			yield put(messagesFailedActionCreator(index)(response.error))
		else
			yield put(updateMessageSucceededActionCreator(index)(messages))
	} catch (e) {
		yield put(messagesFailedActionCreator(index)(e.message))
	}
}

function* deleteMessage({index, payload}) {
	try {
		const response = yield call(api.deleteMessage, payload.ID)

		let messages = yield select(selectMessages(index))
		messages = messages.filter(({ID}) => ID !== payload.ID)

		if (response.error !== "")
			yield put(messagesFailedActionCreator(index)(response.error))
		else
			yield put(deleteMessageSucceededActionCreator(index)(messages))
	} catch (e) {
		yield put(messagesFailedActionCreator(index)(e.message))
	}
}

function* copyMessage({payload}) {
	try {
		const browser = yield select(selectBrowser)
		yield call(async function copy(text, browser) {
			// TODO: compile front with browser dirrectives?
			switch (browser) {
			case "chrome":
				const result = await navigator.permissions.query({ name: "clipboard-write" })
				if (result.state === "granted" || result.state === "prompt")
					await navigator.clipboard.writeText(text)
				else
					console.error("clipboard write permission is not granted")

				break

			case "firefox":
				await navigator.clipboard.writeText(text)
				break
			}
		}, payload.text, browser)
	} catch (e) {
		console.error(e)
	}
}

function* stackSaga() {
	yield takeLatest(DELETE_MESSAGE, deleteMessage)
	yield takeLatest(UPDATE_MESSAGE, updateMessage)
	yield takeLatest(FETCH_MESSAGES, fetchMessages)
	yield takeLatest(SEND_MESSAGE, sendMessage)
	yield takeLatest(COPY_MESSAGE, copyMessage)
}

export {stackSaga}