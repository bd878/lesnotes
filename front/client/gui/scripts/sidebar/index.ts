import onSearchFormSubmit from './onSearchFormSubmit';
import onThemeSettingsClick from './onThemeSettingsClick';
import onLangSettingsClick from './onLangSettingsClick';
import onFontSizeSettingsClick from './onFontSizeSettingsClick';
import getByID from '../getByID'

const elems = {
	form:   document.createElement("form"),
	div:    document.createElement("div"),
	button: document.createElement("button"),
	input:  document.createElement("input"),

	get searchFormElem():        HTMLFormElement     { return getByID("messages-search-form",  this.form) as HTMLFormElement },
	get themeSettingsElem():     HTMLDivElement      { return getByID("theme-settings",        this.div) as HTMLDivElement },
	get langSettingsElem():      HTMLDivElement      { return getByID("lang-settings",         this.div) as HTMLDivElement },
	get fontSizeSettingsElem():  HTMLDivElement      { return getByID("font-size-settings",    this.div) as HTMLDivElement },
}

function init() {
	console.log("sidebar loaded")

	elems.searchFormElem.addEventListener("submit",       e => onSearchFormSubmit(elems, e))
	elems.themeSettingsElem.addEventListener("click",     e => onThemeSettingsClick(elems, e))
	elems.langSettingsElem.addEventListener("click",      e => onLangSettingsClick(elems, e))
	elems.fontSizeSettingsElem.addEventListener("click",  e => onFontSizeSettingsClick(elems, e))
}

window.addEventListener("load", init)
