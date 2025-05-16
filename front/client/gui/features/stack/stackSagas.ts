import {takeLatest,put,call,select} from 'redux-saga/effects'
import {
	PUBLISH_SELECTED,
	PRIVATE_SELECTED,
	UPDATE_MESSAGE,
	FETCH_MESSAGES,
	SEND_MESSAGE,
	DELETE_MESSAGE,
	DELETE_SELECTED,
	COPY_MESSAGE,
	COPY_LINK,
} from './stackActions'
import {
	resetEditMessageActionCreator,
	messagesFailedActionCreator,
	fetchMessagesSucceededActionCreator,
	sendMessageSucceededActionCreator,
	updateMessageSucceededActionCreator,
	deleteSelectedSucceededActionCreator,
} from './stackActionCreators'
import {showNotificationActionCreator} from '../notification';
import * as is from '../../../third_party/is'
import i18n from '../../../i18n';
import {selectMessages, selectSelectedMessageIDs, selectMessageForEdit} from './stackSelectors';
import {selectBrowser, selectIsMobile, selectIsDesktop} from '../me'
import api, {getMessageLinkUrl} from '../../../api'

interface FetchMessagesPayload {
	limit:  number;
	offset: number;
	order:  number;
}

function* fetchMessages({index, payload}: {payload: FetchMessagesPayload}) {
	try {
		const response = yield call(api.loadMessages,
			{limit: payload.limit, offset: payload.offset, order: payload.order, threadID: payload.threadID})

		response.messages.reverse();
		if (is.notEmpty(response.error))
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
			if (is.notEmpty(response.error)) {
				yield put(messagesFailedActionCreator(index)(response.error))
				return
			}

			response = yield call(api.sendMessage, {text: payload.text, fileID: response.ID, threadID: payload.threadID})
		} else {
			response = yield call(api.sendMessage, {text: payload.text, threadID: payload.threadID})
		}

		if (is.notEmpty(response.error))
			yield put(messagesFailedActionCreator(index)(response.error))
		else
			yield put(sendMessageSucceededActionCreator(index)(response.message))
	} catch (e) {
		yield put(messagesFailedActionCreator(index)(e.message))
	}
}

function* updateMessage({index, payload}) {
	try {
		const response = yield call(api.updateMessage, {id: payload.ID, text: payload.text})

		const messages = yield select(selectMessages(index))
		let idx = messages.findIndex(({ID}) => ID === payload.ID)
		if (idx !== -1)
			messages[idx] = {
				...messages[idx],
				ID: response.ID,
				text: payload.text,
				updateUTCNano: response.updateUTCNano,
			}

		if (is.notEmpty(response.error))
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

		const messageForEdit = yield select(selectMessageForEdit(index))

		if (is.notEmpty(response.error)) {
			yield put(messagesFailedActionCreator(index)(response.error))
		} else {
			yield put(deleteMessageSucceededActionCreator(index)(messages))
			if (payload.ID === messageForEdit.ID)
				yield put(resetEditMessageActionCreator(index)({}))
		}
	} catch (e) {
		yield put(messagesFailedActionCreator(index)(e.message))
	}
}

function* deleteSelected({index}) {
	try {
		const idsSet = yield select(selectSelectedMessageIDs(index))
		const response = yield call(api.deleteMessages, Array.from(idsSet))
		let messages = yield select(selectMessages(index))
		messages = messages.filter(({ID}) => !idsSet.has(ID))

		const messageForEdit = yield select(selectMessageForEdit(index))

		if (is.notEmpty(response.error)) {
			yield put(messagesFailedActionCreator(index)(response.error))
		} else {
			yield put(deleteSelectedSucceededActionCreator(index)(messages))
			if (idsSet.has(messageForEdit.ID))
				yield put(resetEditMessageActionCreator(index)({}))
		}
	} catch (e) {
		yield put(messagesFailedActionCreator(index)(e.message))
	}
}

function* copyMessage({payload}) {
	try {
		const browser = yield select(selectBrowser)
		yield call(async function copy(text, browser) {
			// TODO: compile front with browser directives?
			switch (browser) {
			case "chrome":
				const result = await navigator.permissions.query({ name: "clipboard-write" })
				if (result.state === "granted" || result.state === "prompt")
					await navigator.clipboard.writeText(text)
				else
					throw new Error("clipboard write permission is not granted")

				break

			case "firefox":
				await navigator.clipboard.writeText(text)
				break
			}
		}, payload.text, browser)

		yield put(showNotificationActionCreator({text: i18n("copied")}))
	} catch (e) {
		console.error(e)
	}
}

function* copyLink({payload}) {
	try {
		const browser = yield select(selectBrowser)
		yield call(async function copy(id, browser) {
			// TODO: compile front with browser directives?
			const text = getMessageLinkUrl(id)
			switch (browser) {
			case "chrome":
				const result = await navigator.permissions.query({ name: "clipboard-write" })
				if (result.state === "granted" || result.state === "prompt")
					await navigator.clipboard.writeText(text)
				else
					throw new Error("clipboard write permission is not granted")

				break

			case "firefox":
				await navigator.clipboard.writeText(text)
				break
			}
		}, payload.ID, browser)

		yield put(showNotificationActionCreator({text: i18n("copied")}))
	} catch (e) {
		console.error(e)
	}
}

function* publishSelected({index, payload}) {
	try {
		const idsSet = yield select(selectSelectedMessageIDs(index))
		const response = yield call(api.publishMessages, Array.from(idsSet))

		const messages = yield select(selectMessages(index))
		messages.filter(({ID}) => idsSet.has(ID)).forEach(msg => {
			msg.private = false
			msg.updateUTCNano = response.updateUTCNano
		})

		if (is.notEmpty(response.error))
			yield put(messagesFailedActionCreator(index)(response.error))
		else
			yield put(updateMessageSucceededActionCreator(index)(messages))
	} catch (e) {
		yield put(messagesFailedActionCreator(index)(e.message))
	}
}

function* privateSelected({index, payload}) {
	try {
		const idsSet = yield select(selectSelectedMessageIDs(index))
		const response = yield call(api.privateMessages, Array.from(idsSet))

		const messages = yield select(selectMessages(index))
		messages.filter(({ID}) => idsSet.has(ID)).forEach(msg => {
			msg.private = true
			msg.updateUTCNano = response.updateUTCNano
		})

		if (is.notEmpty(response.error))
			yield put(messagesFailedActionCreator(index)(response.error))
		else
			yield put(updateMessageSucceededActionCreator(index)(messages))
	} catch (e) {
		yield put(messagesFailedActionCreator(index)(e.message))
	}
}

function* stackSaga() {
	yield takeLatest(DELETE_MESSAGE, deleteMessage)
	yield takeLatest(DELETE_SELECTED, deleteSelected)
	yield takeLatest(UPDATE_MESSAGE, updateMessage)
	yield takeLatest(FETCH_MESSAGES, fetchMessages)
	yield takeLatest(SEND_MESSAGE, sendMessage)
	yield takeLatest(COPY_MESSAGE, copyMessage)
	yield takeLatest(COPY_LINK, copyLink)
	yield takeLatest(PUBLISH_SELECTED, publishSelected)
	yield takeLatest(PRIVATE_SELECTED, privateSelected)
}

export {stackSaga}