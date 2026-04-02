import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder';

let threadTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/desktop/thread.mustache')), { encoding: 'utf-8' });
let threadTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/mobile/thread.mustache')), { encoding: 'utf-8' });

class PublicThreadBuilder extends AbstractPublicBuilder {
	auth   = undefined
	header = undefined
	tree   = undefined

	addAuth(auth: Builder) {
		this.auth = auth.build()
	}

	addHeader(header: Builder) {
		this.header = header.build()
	}

	addMessagesTree(tree: Builder) {
		this.tree = tree.build()
	}

	build() {
		return mustache.render(this.isMobile ? threadTemplateMobile : threadTemplate, {
		}, {
			auth:  this.auth,
			header: this.header,
			messagesTree:  this.tree,
		})
	}

}

export default PublicThreadBuilder
