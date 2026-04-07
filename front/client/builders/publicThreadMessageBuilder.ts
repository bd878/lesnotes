import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder'

let threadMessageTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/desktop/thread.mustache')), { encoding: 'utf-8' });
let threadMessageTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/mobile/thread.mustache')), { encoding: 'utf-8' });

class PublicThreadMessageBuilder extends AbstractPublicBuilder {
	header             = undefined;
	tree               = undefined;
	messageFeatures    = undefined;
	messageView        = undefined;

	addHeader(header: Builder) {
		this.header = header.build()
		return this
	}

	addMessageFeatures(features: Builder) {
		this.messageFeatures = features.build()
		return this
	}

	addMessagesTree(tree: Builder) {
		this.tree = tree.build()
		return this
	}

	addMessageView(view: Builder) {
		this.messageView = view.build()
		return this
	}

	build() {
		return mustache.render(this.isMobile ? threadMessageTemplateMobile : threadMessageTemplate, {
			hasMessage:        this.messageView != undefined,
		}, {
			header:            this.header,
			messagesTree:      this.tree,
			messageView:       this.messageView,
			messageFeatures:   this.messageFeatures,
		})
	}

}

export default PublicThreadMessageBuilder
