// import autosaveMessage from './autosaveMessage';
import onFormSubmit from './onFormSubmit';
import onFilesListClick from './onFilesListClick';
import debounce from '../../scripts/debounce'
import getByID from '../../scripts/getByID'

const elems = {
	form:      document.createElement("form"),
	input:     document.createElement("input"),
	textarea:  document.createElement("textarea"),
	select:    document.createElement("select"),

	get messageEditFormElem():       HTMLFormElement     { return getByID("edit-message-form",     this.form)   as HTMLFormElement },
	get filesListElem():             HTMLDivElement      { return getByID("files-list",            this.button) as HTMLDivElement },
	get filesInputElem():            HTMLInputElement    { return getByID("files-input",           this.input) as HTMLInputElement },
}

function init() {
	console.log("messageEditScript loaded");

	// elems.messageEditForm.addEventListener("change", debounce(e => autosaveMessage(elems, e), 5000 /* we might save by submit button */))
	elems.messageEditFormElem.addEventListener("submit", e => onFormSubmit(elems, e))
	elems.filesListElem.addEventListener("click",    e => onFilesListClick(elems, e))
}

window.addEventListener("load", init)