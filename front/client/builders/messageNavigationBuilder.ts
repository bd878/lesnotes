import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let messageNavigationTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/message_navigation/desktop/message_navigation.mustache')), { encoding: 'utf-8' });
let messageNavigationTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/message_navigation/mobile/message_navigation.mustache')), { encoding: 'utf-8' });

class MessageNavigationBuilder extends AbstractBuilder {
	attachments = false
	comments = false
	translations = false

	addAttachments() {
		this.attachments = true
	}

	addComments() {
		this.comments = true
	}

	addTranslations() {
		this.translations = true
	}

	build() {
		const search = this.search

		return mustache.render(this.isMobile ? messageNavigationTemplateMobile : messageNavigationTemplate, {
			isAttachments:    this.attachments,
			isComments:       this.comments,
			isTranslations:   this.translations,
			attachments:      this.i18n("attachments"),
			comments:         this.i18n("comments"),
			translations:     this.i18n("translations"),
			attachmentsHref:  function() { const params = new URLSearchParams(search); params.set("nav", "files");     return "?" + params.toString(); },
			commentsHref:     function() { const params = new URLSearchParams(search); params.set("nav", "comments");  return "?" + params.toString(); },
			translationsHref: function() { const params = new URLSearchParams(search); params.set("nav", "trans");     return "?" + params.toString(); },
		}, {})
	}
}

export default MessageNavigationBuilder
