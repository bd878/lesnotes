import type { Message } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder';

let messageViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/message/desktop/message_view.mustache')), { encoding: 'utf-8' });
let messageViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/message/mobile/message_view.mustache')), { encoding: 'utf-8' });

class MessageViewBuilder extends AbstractPublicBuilder {
	message = undefined
	redirectUrl = ""
	deleteRedirectUrl = ""
	editHref = ""

	addRedirectUrl(redirectUrl: string) {
		this.redirectUrl = redirectUrl
		return this
	}

	addDeleteRedirectUrl(deleteRedirectUrl: string) {
		this.deleteRedirectUrl = deleteRedirectUrl
		return this
	}

	addMessage(message: Message) {
		this.message = message
		return this
	}

	addMessageEditHref(href: string) {
		this.editHref = href
		return this
	}

	build() {
		const editHref = this.editHref
		return mustache.render(this.isMobile ? messageViewTemplateMobile : messageViewTemplate, {
			isAuthed:              this.isAuthed,
			message:               this.message,
			editHref:              editHref || `/messages/${this.message.ID}/edit` + this.search,
			deleteAction:          "/message/delete" + this.search,
			publishAction:         "/message/publish" + this.search,
			privateAction:         "/message/private" + this.search,
			redirectUrl:           this.redirectUrl || (this.path + this.search),
			deleteRedirectUrl:     this.deleteRedirectUrl,
			domain:                Config.get("domain"),
		})
	}
}

export default MessageViewBuilder
