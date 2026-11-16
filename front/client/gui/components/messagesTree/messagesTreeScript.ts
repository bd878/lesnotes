import onMessagesListScroll, { restoreMessagesListScroll } from './onMessagesListScroll';
import debounce from '../../scripts/debounce'
import getByID from '../../scripts/getByID'

const elems = {
	div:    document.createElement("div"),

	get messagesListElem():      HTMLDivElement      { return getByID("messages-list",         this.div) as HTMLDivElement },
}

function init() {
	restoreMessagesListScroll(elems)

	elems.messagesListElem.addEventListener('scroll',     debounce(e => onMessagesListScroll(elems, e), 1000))
}

window.addEventListener("load", init)
