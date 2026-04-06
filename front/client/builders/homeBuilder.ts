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

let controlPanelTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/control_panel.mustache')), { encoding: 'utf-8' });
let controlPanelTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/control_panel.mustache')), { encoding: 'utf-8' });

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
	auth                 = undefined;
	messageHeader        = undefined;
	scripts              = ["/public/pages/home/homeScript.js"]

	addMessagesTree(tree: Builder) {
		this.messagesTree = tree.build()
	}

	addMessageHeader(header: Builder) {
		this.messageHeader = header.build()
	}

	addMessageFeatures(features: Builder) {
		this.messageFeatures = features.build()
	}

	addMessageView(view: Builder) {
		this.messageView = view.build()
	}

	addNewMessageForm(form: Builder) {
		this.newMessageForm = form.build()
	}

	addMessageEditForm(form: ScriptsBuilder) {
		this.messageEditForm = form.build()
		this.scripts.push(...form.scripts)
	}

	addControlPanel() {
		this.controlPanel = mustache.render(this.isMobile ? controlPanelTemplateMobile : controlPanelTemplate, {}, {
			logout:           this.auth,
		})
	}

	addAuth(auth: Builder) {
		this.auth = auth.build()
	}

	addHeader(header: Builder) {
		this.header = header.build()
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
