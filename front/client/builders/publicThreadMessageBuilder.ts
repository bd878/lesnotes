import type { Message, TranslationPreview } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder'

let translationsTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/desktop/translations.mustache')), { encoding: 'utf-8' });
let translationsTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/mobile/translations.mustache')), { encoding: 'utf-8' });

let threadMessageTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/desktop/thread.mustache')), { encoding: 'utf-8' });
let threadMessageTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/mobile/thread.mustache')), { encoding: 'utf-8' });

class PublicThreadMessageBuilder extends AbstractPublicBuilder {
	addTranslations(message: number | string, thread: string, previews: TranslationPreview[]) {
		const search = this.search

		this.translations = mustache.render(this.isMobile ? translationsTemplateMobile : translationsTemplate, {
			mainMessage:           this.i18n("mainMessage"),
			mainMessageHref:       function() { return `/t/${thread}/${message}` + search },
			translationHref:       function() { return `/t/${thread}/${message}/${this.lang}` + search },
			translations:          previews,
			hasTranslations:       () => previews.length > 0,
		})
	}

	build(message?: Message) {
		return mustache.render(this.isMobile ? threadMessageTemplateMobile : threadMessageTemplate, {
			message:       message,
		}, {
			signup:            this.signup,
			logout:            this.logout,
			sidebar:           this.sidebar,
			translationView:   this.translationView,
			searchForm:        this.searchForm,
			messagesList:      this.messagesList,
			translations:      this.translations,
			messageView:       this.messageView,
			filesView:         this.filesView,
			comments:          this.comments,
			messageNavigation: this.messageNavigation,
		})
	}

}

export default PublicThreadMessageBuilder
