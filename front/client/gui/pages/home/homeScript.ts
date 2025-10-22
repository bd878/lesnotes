import onFormSubmit from './onFormSubmit';
import onSearchFormSubmit from './onSearchFormSubmit';
import onFileInputChange from './onFileInputChange';
import onSelectFilesClick from './onSelectFilesClick';
import onMessageCancelClick from './onMessageCancelClick';
import onMessageUpdateFormSubmit from './onMessageUpdateFormSubmit';
import onMessagesListClick from './onMessagesListClick';
import onThreadsListClick from './onThreadsListClick';
import onMessageDeleteClick from './onMessageDeleteClick';
import onMessageEditClick from './onMessageEditClick';
import onMessagePublishClick from './onMessagePublishClick';
import onMessagePrivateClick from './onMessagePrivateClick';
import onMessageCancelEditClick from './onMessageCancelEditClick';
import onThemeSettingsClick from './onThemeSettingsClick';
import onLangSettingsClick from './onLangSettingsClick';
import onFontSizeSettingsClick from './onFontSizeSettingsClick';
import {getByID} from '../../../utils'

const elems = {
	form:   document.createElement("form"),
	div:    document.createElement("div"),
	button: document.createElement("button"),
	input:  document.createElement("input"),

	get messageFormElem():       HTMLFormElement     { return getByID("message-form",          this.form) as HTMLFormElement },
	get filesButtonElem():       HTMLButtonElement   { return getByID("select-files-button",   this.button) as HTMLButtonElement },
	get filesListElem():         HTMLDivElement      { return getByID("files-list",            this.button) as HTMLDivElement },
	get filesInputElem():        HTMLInputElement    { return getByID("files-input",           this.input) as HTMLInputElement },
	get editFormElem():          HTMLFormElement     { return getByID("message-edit-form",     this.form) as HTMLFormElement },
	get searchFormElem():        HTMLFormElement     { return getByID("messages-search-form",  this.form) as HTMLFormElement },
	get messagesListElem():      HTMLDivElement      { return getByID("messages-list",         this.div) as HTMLDivElement },
	get threadsListElem():       HTMLDivElement      { return getByID("threads-list",          this.div) as HTMLDivElement },
	get messageDeleteElem():     HTMLButtonElement   { return getByID("message-delete",        this.button) as HTMLButtonElement },
	get messageEditElem():       HTMLButtonElement   { return getByID("message-edit",          this.button) as HTMLButtonElement },
	get messagePublishElem():    HTMLButtonElement   { return getByID("message-publish",       this.button) as HTMLButtonElement },
	get messagePrivateElem():    HTMLButtonElement   { return getByID("message-private",       this.button) as HTMLButtonElement },
	get messageCancelEditElem(): HTMLButtonElement   { return getByID("message-cancel-edit",   this.button) as HTMLButtonElement },
	get messageCancelElem():     HTMLButtonElement   { return getByID("message-cancel",        this.button) as HTMLButtonElement },
	get themeSettingsElem():     HTMLDivElement      { return getByID("theme-settings",        this.div) as HTMLDivElement },
	get langSettingsElem():      HTMLDivElement      { return getByID("lang-settings",         this.div) as HTMLDivElement },
	get fontSizeSettingsElem():  HTMLDivElement      { return getByID("font-size-settings",    this.div) as HTMLDivElement },
}

function init() {
	console.log("loaded")

	// TODO: rewrite on event bus
	elems.messageFormElem.addEventListener("submit",      e => onFormSubmit(elems, e))
	elems.searchFormElem.addEventListener("submit",       e => onSearchFormSubmit(elems, e))
	elems.filesInputElem.addEventListener("change",       e => onFileInputChange(elems, e))
	elems.filesButtonElem.addEventListener("click",       e => onSelectFilesClick(elems, e))
	elems.messageCancelElem.addEventListener("click",     e => onMessageCancelClick(elems, e))
	elems.editFormElem.addEventListener("submit",         e => onMessageUpdateFormSubmit(elems, e))
	elems.messagesListElem.addEventListener("click",      e => onMessagesListClick(elems, e))
	elems.threadsListElem.addEventListener("click",       e => onThreadsListClick(elems, e))
	elems.messageDeleteElem.addEventListener("click",     e => onMessageDeleteClick(elems, e))
	elems.messageEditElem.addEventListener("click",       e => onMessageEditClick(elems, e))
	elems.messagePublishElem.addEventListener("click",    e => onMessagePublishClick(elems, e))
	elems.messagePrivateElem.addEventListener("click",    e => onMessagePrivateClick(elems, e))
	elems.messageCancelEditElem.addEventListener("click", e => onMessageCancelEditClick(elems, e))
	elems.themeSettingsElem.addEventListener("click",     e => onThemeSettingsClick(elems, e))
	elems.langSettingsElem.addEventListener("click",      e => onLangSettingsClick(elems, e))
	elems.fontSizeSettingsElem.addEventListener("click",  e => onFontSizeSettingsClick(elems, e))
}

window.addEventListener("load", init)
