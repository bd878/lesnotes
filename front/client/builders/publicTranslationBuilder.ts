import type { Translation, TranslationPreview } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder'

class PublicTranslationBuilder extends AbstractPublicBuilder {
	async addTranslations(message: number | string, previews: TranslationPreview[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/translations/mobile/translations.mustache' : 'templates/translations/desktop/translations.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.translations = mustache.render(template, {
			mainMessage:           this.i18n("mainMessage"),
			mainMessageHref:       function() { return `/m/${message}` },
			translationHref:       function() { return `/m/${message}/${this.lang}` },
			translations:          previews,
			hasTranslations:       () => previews.length > 0,
		})
	}

	async build(translation?: Translation) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/layout/mobile/layout.mustache' : 'templates/layout/desktop/layout.mustache'
		)), { encoding: 'utf-8' });
		const content = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/translation/mobile/translation.mustache' : 'templates/translation/desktop/translation.mustache'
		)), { encoding: 'utf-8' });

		const theme = this.theme
		const fontSize = this.fontSize

		return mustache.render(layout, {
			html: () => (text, render) => {
				let html = "<html"

				if (theme) html += ` class="${theme}"`;
				if (this.lang) html += ` lang="${this.lang}"`;
				if (fontSize) html += ` data-size="${fontSize}"`
				html += ">"

				return html + render(text) + "</html>"
			},
			manifest: "/public/manifest.json",
			styles:   styles,
			lang:     this.lang,
			theme:    this.theme,
		}, {
			footer: this.footer,
			content: mustache.render(content, {
				translation:      translation,
			}, {
				signup:           this.signup,
				logout:           this.logout,
				sidebar:          this.sidebar,
				translationView:  this.translationView,
				translations:     this.translations,
				filesView:        this.filesView,
			})
		})
	}

}

export default PublicTranslationBuilder