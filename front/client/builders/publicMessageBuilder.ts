import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder'

let messageTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/message/desktop/message.mustache')), { encoding: 'utf-8' });
let messageTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/message/mobile/message.mustache')), { encoding: 'utf-8' });

class PublicMessageBuilder extends AbstractPublicBuilder {
	messageFeatures    = undefined;
	messageView        = undefined;
	header             = undefined;

	addMessageFeatures(features: Builder) {
		this.messageFeatures = features.build()
	}

	addMessageView(message: Builder) {
		this.messageView = message.build()
	}

	addHeader(header: Builder) {
		this.header = header.build()
	}

	build() {
		return mustache.render(this.isMobile ? messageTemplateMobile : messageTemplate, {}, {
			messageView:       this.messageView,
			messageFeatures:   this.messageFeatures,
			header:            this.header,
		})
	}

}

export default PublicMessageBuilder
