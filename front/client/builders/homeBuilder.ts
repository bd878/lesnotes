import type {Builder, ScriptsBuilder} from './builder';
import type {File, Message} from '../api/models';
import type {FileWithMime} from '../types';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import { unwrapPaging } from '../api/models/paging';
import * as is from '../third_party/is';
import i18n from '../i18n';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'

let homeTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/home.mustache')), { encoding: 'utf-8' });
let homeTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/home.mustache')), { encoding: 'utf-8' });

class HomeBuilder extends AbstractBuilder {
	messageEditForm      = undefined;
	messageView          = undefined;
	newMessageForm       = undefined;
	threadView           = undefined;
	threadEditForm       = undefined;
	header               = undefined;
	messagesTree         = undefined;
	messageFeatures      = undefined;
	controlPanel         = undefined;
	messageHeader        = undefined;
	scripts              = ["/public/pages/home/homeScript.js"]

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

	addThreadEditForm(form: ScriptsBuilder) {
		this.threadEditForm = form.build()
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

	build() {
		return mustache.render(this.isMobile ? homeTemplateMobile : homeTemplate, {
			hasFeatures:          this.messageFeatures != undefined,
		}, {
			messageEditForm:      this.messageEditForm,
			messageView:          this.messageView,
			threadView:           this.threadView,
			header:               this.header,
			messageFeatures:      this.messageFeatures,
			threadEditForm:       this.threadEditForm,
			newMessageForm:       this.newMessageForm,
			messagesTree:         this.messagesTree,
			controlPanel:         this.controlPanel,
			messageHeader:        this.messageHeader,
		});
	}
}

export default HomeBuilder
