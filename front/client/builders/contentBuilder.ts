import type {Builder, ScriptsBuilder} from './builder';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'

let contentTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/content/desktop/content.mustache')), { encoding: 'utf-8' });
let contentTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/content/mobile/content.mustache')), { encoding: 'utf-8' });

class ContentBuilder extends AbstractBuilder {
	messageEditForm      = undefined;
	messageView          = undefined;
	newMessageForm       = undefined;
	threadView           = undefined;
	threadPath           = undefined;
	header               = undefined;
	messagesTree         = undefined;
	messageFeatures      = undefined;
	controlPanel         = undefined;
	messageHeader        = undefined;

	addMessagesTree(tree: Builder) {
		this.messagesTree = tree.build()
		return this
	}

	addMessageHeader(header: Builder) {
		this.messageHeader = header.build()
		return this
	}

	addMessageFeatures(features: Builder) {
		this.messageFeatures = features.build()
		return this
	}

	addMessageView(view: Builder) {
		this.messageView = view.build()
		return this
	}

	addNewMessageForm(form: Builder) {
		this.newMessageForm = form.build()
		return this
	}

	addMessageEditForm(form: ScriptsBuilder) {
		this.messageEditForm = form.build()
		this.scripts.push(...form.scripts)
		return this
	}

	addThreadView(view: ScriptsBuilder) {
		this.threadView = view.build()
		this.scripts.push(...view.scripts)
		return this
	}

	addControlPanel(panel: Builder) {
		this.controlPanel = panel.build()
		return this
	}

	addHeader(header: Builder) {
		this.header = header.build()
		return this
	}

	addThreadPath(threadPath: Builder) {
		this.threadPath = threadPath.build()
		return this
	}

	build() {
		return mustache.render(this.isMobile ? contentTemplateMobile : contentTemplate, {
			hasFeatures:          this.messageFeatures != undefined,
			hasMessage:           this.messageView != undefined
						|| this.messageEditForm != undefined
						|| this.newMessageForm != undefined,
			hasMessagesTree:      this.messagesTree != undefined,
		}, {
			threadPath:           this.threadPath,
			messageEditForm:      this.messageEditForm,
			messageView:          this.messageView,
			threadView:           this.threadView,
			header:               this.header,
			messageFeatures:      this.messageFeatures,
			newMessageForm:       this.newMessageForm,
			messagesTree:         this.messagesTree,
			controlPanel:         this.controlPanel,
			messageHeader:        this.messageHeader,
		});
	}
}

export default ContentBuilder
