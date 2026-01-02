import onMessagesListDragStart from './onMessagesListDragStart';
import onMessagesListDrop from './onMessagesListDrop';
import getByID from '../../scripts/getByID'

const elems = {
	form:   document.createElement("form"),
	div:    document.createElement("div"),
	button: document.createElement("button"),
	input:  document.createElement("input"),

	get messagesListElem():      HTMLDivElement      { return getByID("messages-list",         this.div) as HTMLDivElement },
}

function init() {
	console.log("loaded")

	elems.messagesListElem.addEventListener('dragstart',  e => onMessagesListDragStart(elems, e))
	elems.messagesListElem.addEventListener('drop',       e => onMessagesListDrop(elems, e))
}

window.addEventListener("load", init)
