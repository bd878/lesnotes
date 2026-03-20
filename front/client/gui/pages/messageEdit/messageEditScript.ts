import autosaveMessage from './autosaveMessage';
import debounce from '../../scripts/debounce'
import getByID from '../../scripts/getByID'

const elems = {
	form:      document.createElement("form"),
	input:     document.createElement("input"),
	textarea:  document.createElement("textarea"),
	select:    document.createElement("select"),

	get messageEditForm(): HTMLFormElement { return getByID("message-edit-form", this.form) as HTMLFormElement },
}

function init() {
	console.log("messageEditScript loaded");

	elems.messageEditForm.addEventListener("change", debounce(e => autosaveMessage(elems, e), 5000 /* we might save by submit button */))
}

window.addEventListener("load", init)