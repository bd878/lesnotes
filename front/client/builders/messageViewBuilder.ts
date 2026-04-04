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
	messageView = undefined
	redirectUrl = ""
	deleteRedirectUrl = ""

	addRedirectUrl(redirectUrl: string) {
		this.redirectUrl = redirectUrl
	}

	addDeleteRedirectUrl(deleteRedirectUrl: string) {
		this.deleteRedirectUrl = deleteRedirectUrl
	}

	addMessage(message: Message) {
		this.messageView = mustache.render(this.isMobile ? messageViewTemplateMobile : messageViewTemplate, {
			isAuthed:              this.isAuthed,
			message:               message,
			editHref:              `/editor/messages/${message.ID}` + this.search,
			deleteAction:          "/m/delete" + this.search,
			publishAction:         "/m/publish" + this.search,
			privateAction:         "/m/private" + this.search,
			redirectUrl:           this.redirectUrl || (this.path + this.search),
			deleteRedirectUrl:     this.deleteRedirectUrl,
			domain:                Config.get("domain"),
		})
	}

	build() {
		return this.messageView
	}
}

export default MessageViewBuilder
