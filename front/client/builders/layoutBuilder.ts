import type {Builder, ScriptsBuilder} from './builder';
import Config from 'config';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let footerTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/footer/desktop/footer.mustache')), { encoding: 'utf-8' });
let footerTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/footer/mobile/footer.mustache')), { encoding: 'utf-8' });

let layoutTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/layout/desktop/layout.mustache')), { encoding: 'utf-8' });
let layoutTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/layout/mobile/layout.mustache')), { encoding: 'utf-8' });

let stylesTemplate = readFileSync(resolve(join(Config.get('basedir'),'public/styles/styles.css')), { encoding: 'utf-8' });

class LayoutBuilder extends AbstractBuilder {
	footer        = undefined;
	header        = undefined;
	content       = undefined;

	addHeader(header: Builder) {
		this.header = header.build()
		return this
	}

	addFooter() {
		this.footer = mustache.render(this.isMobile ? footerTemplateMobile : footerTemplate, {
			terms:            this.i18n("terms"),
			contact:          this.i18n("contact"),
			docs:             this.i18n("docs"),
			termsHref:        "/terms" + this.search,
			contactHref:      "/contact" + this.search,
			docsHref:         "/docs" + this.search,
		})
		return this
	}

	addContent(content: ScriptsBuilder) {
		this.content = content.build()
		this.scripts = content.scripts
		return this
	}

	build() {
		const lang = this.lang
		const theme = this.theme
		const fontSize = this.fontSize

		const header = this.header ? this.header : ""
		const content = this.content ? this.content : ""

		return mustache.render(this.isMobile ? layoutTemplateMobile : layoutTemplate, {
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
