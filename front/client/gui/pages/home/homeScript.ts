import onFormSubmit from './onFormSubmit';
import onMessagesListDragStart from './onMessagesListDragStart';
import onMessagesListDrop from './onMessagesListDrop';
import onFileInputChange from './onFileInputChange';
import onSelectFilesClick from './onSelectFilesClick';
import getByID from '../../scripts/getByID'

const elems = {
	form:   document.createElement("form"),
	div:    document.createElement("div"),
	button: document.createElement("button"),
	input:  document.createElement("input"),

	get newMessageFormElem():    HTMLFormElement     { return getByID("new-message-form",      this.form) as HTMLFormElement },
	get messagesListElem():      HTMLDivElement      { return getByID("messages-list",         this.div) as HTMLDivElement },
	get filesInputElem():        HTMLInputElement    { return getByID("files-input",           this.input) as HTMLInputElement },
	get filesListElem():         HTMLDivElement      { return getByID("files-list",            this.button) as HTMLDivElement },
	get filesButtonElem():       HTMLButtonElement   { return getByID("select-files-button",   this.button) as HTMLButtonElement },
}

function init() {
	console.log("homeScript loaded")

	elems.newMessageFormElem.addEventListener("submit",   e => onFormSubmit(elems, e))
	elems.messagesListElem.addEventListener('dragstart',  e => onMessagesListDragStart(elems, e))
	elems.filesButtonElem.addEventListener("click",       e => onSelectFilesClick(elems, e))
	elems.messagesListElem.addEventListener('drop',       e => onMessagesListDrop(elems, e))
	elems.filesInputElem.addEventListener("change",       e => onFileInputChange(elems, e))
}

window.addEventListener("load", init)
