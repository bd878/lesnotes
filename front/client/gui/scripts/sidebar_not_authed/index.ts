import onThemeSettingsClick from './onThemeSettingsClick';
import onLangSettingsClick from './onLangSettingsClick';
import onFontSizeSettingsClick from './onFontSizeSettingsClick';
import onEnterPress from '../onEnterPress';
import getByID from '../getByID'

const elems = {
	form:   document.createElement("form"),
	div:    document.createElement("div"),

	get themeSettingsElem():     HTMLDivElement      { return getByID("theme-settings",        this.div) as HTMLDivElement },
	get langSettingsElem():      HTMLDivElement      { return getByID("lang-settings",         this.div) as HTMLDivElement },
	get fontSizeSettingsElem():  HTMLDivElement      { return getByID("font-size-settings",    this.div) as HTMLDivElement },
}

function init() {
	console.log("sidebar loaded");

	elems.themeSettingsElem.addEventListener("click",        e => onThemeSettingsClick(elems, e))
	elems.langSettingsElem.addEventListener("click",         e => onLangSettingsClick(elems, e))
	elems.fontSizeSettingsElem.addEventListener("click",     e => onFontSizeSettingsClick(elems, e))
	elems.themeSettingsElem.addEventListener("keypress",     onEnterPress(e => onThemeSettingsClick(elems, e)))
	elems.langSettingsElem.addEventListener("keypress",      onEnterPress(e => onLangSettingsClick(elems, e)))
	elems.fontSizeSettingsElem.addEventListener("keypress",  onEnterPress(e => onFontSizeSettingsClick(elems, e)))
}

window.addEventListener("load", init)