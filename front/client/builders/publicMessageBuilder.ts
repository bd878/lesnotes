import type {Builder} from './builder'
import type { Message, TranslationPreview } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder'

let messageTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/message/desktop/message.mustache')), { encoding: 'utf-8' });
let messageTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/message/mobile/message.mustache')), { encoding: 'utf-8' });

let messageViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/message/desktop/message_view.mustache')), { encoding: 'utf-8' });
let messageViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/message/mobile/message_view.mustache')), { encoding: 'utf-8' });

class PublicMessageBuilder extends AbstractPublicBuilder {
	auth               = undefined;
	messageFeatures    = undefined;
	messageView        = undefined;
	header             = undefined;

	addAuth(auth: Builder) {
		this.auth = auth.build()
	}

	addMessageFeatures(features: Builder) {
		this.messageFeatures = features.build()
	}

	addMessageView(message: Message) {
		this.messageView = mustache.render(this.isMobile ? messageViewTemplateMobile : messageViewTemplate, {
			isAuthed:              this.isAuthed,
			message:               message,
			editHref:              `/editor/messages/${message.ID}` + this.search,
			deleteAction:          "/m/delete" + this.search,
			publishAction:         "/m/publish" + this.search,
			privateAction:         "/m/private" + this.search,
		})
	}

	addHeader(header: Builder) {
		this.header = header.build()
	}

	build() {
		return mustache.render(this.isMobile ? messageTemplateMobile : messageTemplate, {}, {
			auth:              this.auth,
			messageView:       this.messageView,
			messageFeatures:   this.messageFeatures,
			header:            this.header,
		})
	}

}

export default PublicMessageBuilder
