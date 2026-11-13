import * as is from '../../../third_party/is';

const scrollKey = "messages-list-scroll"

function onMessagesListScroll(elems, _e) {
	window.sessionStorage.setItem(scrollKey, elems.messagesListElem.scrollTop)
}

function restoreMessagesListScroll(elems) {
	const savedPos = window.sessionStorage.getItem(scrollKey)
	if (is.notEmpty(savedPos)) {
		elems.messagesListElem.scrollTo({ top: parseInt(savedPos, 0), behavior: "instant" })
	}
}

export default onMessagesListScroll
export { restoreMessagesListScroll }
