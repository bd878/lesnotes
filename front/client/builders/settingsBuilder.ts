import Config from 'config';
import * as is from '../third_party/is';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let settingsTemplate = readFileSync(resolve(join(Config.get("basedir"), 'templates/settings/desktop/settings.mustache')), { encoding: 'utf-8' })
let settingsTemplateMobile = readFileSync(resolve(join(Config.get("basedir"), 'templates/settings/mobile/settings.mustache')), { encoding: 'utf-8' })

class SettingsBuilder extends AbstractBuilder {
	build() {
		const search = this.search
		const theme = this.theme
		const fontSize = this.fontSize
		const lang = this.lang

		return mustache.render(this.isMobile ? settingsTemplateMobile : settingsTemplate, {
			settingsHeader:  this.i18n("settingsHeader"),
			fontSizeHeader:  this.i18n("fontSizeHeader"),
			updateButton:    this.i18n("updateButton"),
			langHeader:      this.i18n("langHeader"),
			themeHeader:     this.i18n("themeHeader"),
			themes:          [{theme: "dark", label: this.i18n("darkTheme")}, {theme: "light", label: this.i18n("lightTheme")}],
			fonts:           [{font: "small", label: "aA", css: "text-sm"}, {font: "medium", label: "aA", css: "text-lg"}, {font: "large", label: "aA", css: "text-xl"}],
			langs:           [{lang: "de", label: this.i18n("deLang")}, {lang: "en", label: this.i18n("enLang")}, {lang: "fr", label: this.i18n("frLang")}, {lang: "ru", label: this.i18n("ruLang")}],
			myTheme:         function() { return this.theme == theme },
			myLang:          function() { return this.lang == lang },
			myFont:          function() { return is.notEmpty(fontSize) ? this.font == fontSize : false },
			font:            fontSize,
			theme:           theme,
			lang:            lang,
			fontHref:        function() { const params = new URLSearchParams(search); params.set("size",  this.font);  return "?" + params.toString(); },
			themeHref:       function() { const params = new URLSearchParams(search); params.set("theme", this.theme); return "?" + params.toString(); },
			langHref:        function() { const params = new URLSearchParams(search); params.set("lang",  this.lang);  return "?" + params.toString(); },
		})
	}
}

export default SettingsBuilder
