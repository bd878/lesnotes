import onFormSubmit from './onFormSubmit';
import onFileInputChange from './onFileInputChange';
import onSelectFilesClick from './onSelectFilesClick';
import onMessageCancelClick from './onMessageCancelClick';
import onMessageUpdateFormSubmit from './onMessageUpdateFormSubmit';
import onMessageDeleteClick from './onMessageDeleteClick';
import onMessageEditClick from './onMessageEditClick';
import onMessagePublishClick from './onMessagePublishClick';
import onMessagePrivateClick from './onMessagePrivateClick';
import onMessageCancelEditClick from './onMessageCancelEditClick';
import onMessagesListDragStart from './onMessagesListDragStart';
import onMessagesListDragOver from './onMessagesListDragOver';
import onMessagesListDrop from './onMessagesListDrop';
import getByID from '../../scripts/getByID'

const elems = {
	form:   document.createElement("form"),
	div:    document.createElement("div"),
	button: document.createElement("button"),
	input:  document.createElement("input"),

	get messageFormElem():       HTMLFormElement     { return getByID("message-form",          this.form) as HTMLFormElement },
	get filesButtonElem():       HTMLButtonElement   { return getByID("select-files-button",   this.button) as HTMLButtonElement },
	get filesListElem():         HTMLDivElement      { return getByID("files-list",            this.button) as HTMLDivElement },
	get noFilesElem():           HTMLDivElement      { return getByID("no-files",              this.div) as HTMLDivElement },
	get filesInputElem():        HTMLInputElement    { return getByID("files-input",           this.input) as HTMLInputElement },
	get editFormElem():          HTMLFormElement     { return getByID("message-edit-form",     this.form) as HTMLFormElement },
	get messagesListElem():      HTMLDivElement      { return getByID("messages-list",         this.div) as HTMLDivElement },
	get threadsListElem():       HTMLDivElement      { return getByID("threads-list",          this.div) as HTMLDivElement },
	get messageDeleteElem():     HTMLButtonElement   { return getByID("message-delete",        this.button) as HTMLButtonElement },
	get messageEditElem():       HTMLButtonElement   { return getByID("message-edit",          this.button) as HTMLButtonElement },
	get messagePublishElem():    HTMLButtonElement   { return getByID("message-publish",       this.button) as HTMLButtonElement },
	get messagePrivateElem():    HTMLButtonElement   { return getByID("message-private",       this.button) as HTMLButtonElement },
	get messageCancelEditElem(): HTMLButtonElement   { return getByID("message-cancel-edit",   this.button) as HTMLButtonElement },
	get messageCancelElem():     HTMLButtonElement   { return getByID("message-cancel",        this.button) as HTMLButtonElement }
}

function init() {
	console.log("loaded")

	// TODO: rewrite on event bus
	elems.messageFormElem.addEventListener("submit",      e => onFormSubmit(elems, e))
	elems.filesInputElem.addEventListener("change",       e => onFileInputChange(elems, e))
	elems.filesButtonElem.addEventListener("click",       e => onSelectFilesClick(elems, e))
	elems.messageCancelElem.addEventListener("click",     e => onMessageCancelClick(elems, e))
	elems.editFormElem.addEventListener("submit",         e => onMessageUpdateFormSubmit(elems, e))
	elems.messagesListElem.addEventListener('dragstart',  e => onMessagesListDragStart(elems, e))
	elems.messagesListElem.addEventListener('dragover',   e => onMessagesListDragOver(elems, e))
	elems.messagesListElem.addEventListener('drop',       e => onMessagesListDrop(elems, e))
	elems.messageDeleteElem.addEventListener("click",     e => onMessageDeleteClick(elems, e))
	elems.messageEditElem.addEventListener("click",       e => onMessageEditClick(elems, e))
	elems.messagePublishElem.addEventListener("click",    e => onMessagePublishClick(elems, e))
	elems.messagePrivateElem.addEventListener("click",    e => onMessagePrivateClick(elems, e))
	elems.messageCancelEditElem.addEventListener("click", e => onMessageCancelEditClick(elems, e))
}

window.addEventListener("load", init)
