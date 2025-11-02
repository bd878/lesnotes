import createTgAuth from '../../scripts/createTgAuth';
import onThemeSettingsClick from './onThemeSettingsClick';
import onLangSettingsClick from './onLangSettingsClick';
import onFontSizeSettingsClick from './onFontSizeSettingsClick';
import {getByID} from '../../../utils'

const elems = {
	div:    document.createElement("div"),

	get widgetElem():            HTMLDivElement      { return getByID("telegram-login-widget", this.div) as HTMLDivElement },
	get themeSettingsElem():     HTMLDivElement      { return getByID("theme-settings",        this.div) as HTMLDivElement },
	get langSettingsElem():      HTMLDivElement      { return getByID("lang-settings",         this.div) as HTMLDivElement },
	get fontSizeSettingsElem():  HTMLDivElement      { return getByID("font-size-settings",    this.div) as HTMLDivElement },
}

function init() {
	console.log("loaded")

	elems.widgetElem.appendChild(createTgAuth())

	elems.themeSettingsElem.addEventListener("click",     e => onThemeSettingsClick(elems, e))
	elems.langSettingsElem.addEventListener("click",      e => onLangSettingsClick(elems, e))
	elems.fontSizeSettingsElem.addEventListener("click",  e => onFontSizeSettingsClick(elems, e))
}

window.addEventListener("load", init)
