import type { Message, TranslationPreview } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder'

let translationsTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/desktop/translations.mustache')), { encoding: 'utf-8' });
let translationsTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/mobile/translations.mustache')), { encoding: 'utf-8' });

let messageTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/message/desktop/message.mustache')), { encoding: 'utf-8' });
let messageTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/message/mobile/message.mustache')), { encoding: 'utf-8' });

class PublicMessageBuilder extends AbstractPublicBuilder {
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
		return mustache.render(this.isMobile ? messageTemplateMobile : messageTemplate, {}, {
			signup:            this.signup,
			logout:            this.logout,
			sidebar:           this.sidebar,
			messagesList:      this.messagesList,
			translations:      this.translations,
			messageView:       this.messageView,
			filesView:         this.filesView,
			comments:          this.comments,
			messageNavigation: this.messageNavigation,
		})
	}

}

export default PublicMessageBuilder
