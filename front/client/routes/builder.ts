import Config from 'config';
import i18n from '../i18n';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

abstract class Builder {
	isMobile:      boolean = false;
	lang:          string  = "en";
	search:        string = "";
	path:          string = "";
	theme:         string = "";
	fontSize:      string = "";

	constructor(isMobile: boolean, lang: string = "en", theme: string = "light", fontSize: string = "medium", search: string = "", path: string = "") {
		this.search = search
		this.isMobile = isMobile
		this.lang = lang
		this.path = path
		this.theme = theme
		this.fontSize = fontSize
	}

	i18n(key: string): string {
		return i18n(this.lang)(key)
	}

	abstract build();

	footer = undefined;
	async addFooter() {
		const template = await readFile(resolve(join(Config.get("basedir"), 'templates/footer.mustache')), { encoding: 'utf-8' });

		this.footer = mustache.render(template, {
			terms:            this.i18n("terms"),
			contact:          this.i18n("contact"),
			docs:             this.i18n("docs"),
		})
	}

	settings = undefined;
	async addSettings() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/settings/mobile/settings.mustache' : 'templates/settings/desktop/settings.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search
		const theme = this.theme
		const fontSize = this.fontSize
		const lang = this.lang

		this.settings = mustache.render(template, {
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

export default Builder
