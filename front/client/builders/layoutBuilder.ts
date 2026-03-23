import type {Builder} from './builder';
import Config from 'config';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let footerTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/footer/desktop/footer.mustache')), { encoding: 'utf-8' });
let footerTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/footer/mobile/footer.mustache')), { encoding: 'utf-8' });

let settingsTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/settings/desktop/settings.mustache')), { encoding: 'utf-8' });
let settingsTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/settings/mobile/settings.mustache')), { encoding: 'utf-8' });

let layoutTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/layout/desktop/layout.mustache')), { encoding: 'utf-8' });
let layoutTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/layout/mobile/layout.mustache')), { encoding: 'utf-8' });

let stylesTemplate = readFileSync(resolve(join(Config.get('basedir'),'public/styles/styles.css')), { encoding: 'utf-8' });

class LayoutBuilder extends AbstractBuilder {
	footer        = undefined;
	header        = undefined;
	settings      = undefined;
	content       = undefined;

	scripts       = [];

	addHeader(header: Builder) {
		this.header = header.build()
	}

	addFooter() {
		this.footer = mustache.render(this.isMobile ? footerTemplate : footerTemplateMobile, {
			terms:            this.i18n("terms"),
			contact:          this.i18n("contact"),
			docs:             this.i18n("docs"),
		})
	}

	addContent(content: Builder) {
		this.content = content.build()
	}

	addSettings() {
		const search = this.search
		const theme = this.theme
		const fontSize = this.fontSize
		const lang = this.lang

		this.settings = mustache.render(this.isMobile ? settingsTemplate : settingsTemplateMobile, {
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

	build() {
		const lang = this.lang
		const theme = this.theme
		const fontSize = this.fontSize

		const header = this.header ? this.header : ""
		const content = this.content ? this.content : ""

		return mustache.render(this.isMobile ? layoutTemplate : layoutTemplateMobile, {
			html: () => (text, render) => {
				let html = "<html"

				if (theme) html += ` class="${theme}"`;
				if (lang) html += ` lang="${lang}"`;
				if (fontSize) html += ` data-size="${fontSize}"`
				html += ">"

				return html + render(text) + "</html>"
			},
			scripts:  this.scripts,
			manifest: "/public/manifest.json",
			styles:   stylesTemplate,
			lang:     lang,
			theme:    theme,
		}, {
			footer: this.footer,
			header: header,
			content: content,
		});
	}
}

export default LayoutBuilder
