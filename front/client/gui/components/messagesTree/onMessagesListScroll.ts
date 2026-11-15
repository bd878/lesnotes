import * as is from '../../../third_party/is';

const scrollKey = "messages-list-scroll"

function onMessagesListScroll(elems, _e) {
	const threadID = getThreadID(elems.messagesListElem)
	const storageScrollJson = setStorageScrollPosition(threadID, elems.messagesListElem.scrollTop)
	window.sessionStorage.setItem(scrollKey, storageScrollJson)
}

function restoreMessagesListScroll(elems) {
	const threadID = getThreadID(elems.messagesListElem)
	const savedPos = getStorageScrollPosition(threadID)
	elems.messagesListElem.scrollTo({ top: savedPos, behavior: "instant" })
}

function getThreadID(treeElem): string {
	const threadIDStr = treeElem.dataset.threadId
	if (is.empty(threadIDStr)) {
		return "0"
	}
	return threadIDStr
}

function getStorageScrollPosition(key: string): number {
	const scrollMap = getScrollPosMap()
	if (scrollMap.has(key)) {
		return parseInt(scrollMap.get(key), 0)
	}
	return 0
}

function setStorageScrollPosition(key: string, scrollTop: number): string {
	let scrollMap = getScrollPosMap()
	scrollMap.set(key, scrollTop)
	return serializeMap(scrollMap)
}

function getScrollPosMap() {
	const posJson = window.sessionStorage.getItem(scrollKey)
	if (is.empty(posJson)) {
		return new Map()
	}
	return new Map(Object.entries(JSON.parse(posJson)))
}

function serializeMap(map): string {
	return JSON.stringify(Object.fromEntries(map))
}

export default onMessagesListScroll
export { restoreMessagesListScroll }
