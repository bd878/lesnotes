import type {Builder} from './builder'
import type { Message, TranslationPreview } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder'

let threadMessageTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/desktop/thread.mustache')), { encoding: 'utf-8' });
let threadMessageTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/mobile/thread.mustache')), { encoding: 'utf-8' });

class PublicThreadMessageBuilder extends AbstractPublicBuilder {
	auth               = undefined;
	messageFeatures    = undefined;
	messageView        = undefined;

	addAuth(auth: Builder) {
		this.auth = auth.build()
	}

	addMessageFeatures(features: Builder) {
		this.messageFeatures = features.build()
	}

	build(message?: Message) {
		return mustache.render(this.isMobile ? threadMessageTemplateMobile : threadMessageTemplate, {
			message:       message,
		}, {
			auth:              this.auth,
			messageView:       this.messageView,
			messageFeatures:   this.messageFeatures,
		})
	}

}

export default PublicThreadMessageBuilder
