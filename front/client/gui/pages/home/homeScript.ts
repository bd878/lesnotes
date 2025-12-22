import onFileInputChange from './onFileInputChange';
import onSelectFilesClick from './onSelectFilesClick';
import onMessagesListDragStart from './onMessagesListDragStart';
import onMessagesListDrop from './onMessagesListDrop';
import getByID from '../../scripts/getByID'

const elems = {
	form:   document.createElement("form"),
	div:    document.createElement("div"),
	button: document.createElement("button"),
	input:  document.createElement("input"),

	get filesButtonElem():       HTMLButtonElement   { return getByID("select-files-button",   this.button) as HTMLButtonElement },
	get filesListElem():         HTMLDivElement      { return getByID("files-list",            this.button) as HTMLDivElement },
	get noFilesElem():           HTMLDivElement      { return getByID("no-files",              this.div) as HTMLDivElement },
	get filesInputElem():        HTMLInputElement    { return getByID("files-input",           this.input) as HTMLInputElement },
	get messagesListElem():      HTMLDivElement      { return getByID("messages-list",         this.div) as HTMLDivElement },
}

function init() {
	console.log("loaded")

	elems.filesInputElem.addEventListener("change",       e => onFileInputChange(elems, e))
	elems.filesButtonElem.addEventListener("click",       e => onSelectFilesClick(elems, e))
	elems.messagesListElem.addEventListener('dragstart',  e => onMessagesListDragStart(elems, e))
	elems.messagesListElem.addEventListener('drop',       e => onMessagesListDrop(elems, e))
}

window.addEventListener("load", init)
