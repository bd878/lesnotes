import type { Translation, TranslationPreview } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder'

let translationsTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/desktop/translations.mustache')), { encoding: 'utf-8' });
let translationsTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/mobile/translations.mustache')), { encoding: 'utf-8' });

let translationTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translation/desktop/translation.mustache')), { encoding: 'utf-8' });
let translationTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translation/mobile/translation.mustache')), { encoding: 'utf-8' });

class PublicTranslationBuilder extends AbstractPublicBuilder {
	addTranslations(message: number | string, previews: TranslationPreview[]) {
		const search = this.search

		this.translations = mustache.render(this.isMobile ? translationsTemplateMobile : translationsTemplate, {
			mainMessage:           this.i18n("mainMessage"),
			mainMessageHref:       function() { return `/m/${message}` + search },
			translationHref:       function() { return `/m/${message}/${this.lang}` + search },
			translations:          previews,
			hasTranslations:       () => previews.length > 0,
		})
	}

	build() {
		return mustache.render(this.isMobile ? translationTemplateMobile : translationTemplate, {
		}, {
			signup:            this.signup,
			logout:            this.logout,
			sidebar:           this.sidebar,
			translationView:   this.translationView,
			translations:      this.translations,
			filesView:         this.filesView,
			comments:          this.comments,
			messageNavigation: this.messageNavigation,
		})
	}

}

export default PublicTranslationBuilder