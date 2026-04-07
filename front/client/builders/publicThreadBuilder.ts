import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder';

let threadTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/desktop/thread.mustache')), { encoding: 'utf-8' });
let threadTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/mobile/thread.mustache')), { encoding: 'utf-8' });

class PublicThreadBuilder extends AbstractPublicBuilder {
	header        = undefined
	tree          = undefined
	controlPanel  = undefined

	addHeader(header: Builder) {
		this.header = header.build()
		return this
	}

	addMessagesTree(tree: Builder) {
		this.tree = tree.build()
		return this
	}

	addControlPanel(panel: Builder) {
		this.controlPanel = panel.build()
		return this
	}

	build() {
		return mustache.render(this.isMobile ? threadTemplateMobile : threadTemplate, {
		}, {
			header:        this.header,
			messagesTree:  this.tree,
			controlPanel:  this.controlPanel
		})
	}

}

export default PublicThreadBuilder
