import type {Message} from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder';

let messagesTreeTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/messages_tree.mustache')), { encoding: 'utf-8' });
let messagesTreeTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/messages_tree.mustache')), { encoding: 'utf-8' });

class MessagesTreeBuilder extends AbstractPublicBuilder {
	list = undefined
	threadPath = undefined

	addList(list: Builder) {
		this.list = list.build()
		return this
	}

	addThreadPath(threadPath: Builder) {
		this.threadPath = threadPath.build()
		return this
	}

	build() {
		return mustache.render(this.isMobile ? messagesTreeTemplateMobile : messagesTreeTemplate, {}, {
			list:       this.list,
			threadPath: this.threadPath,
		})
	}
}

export default MessagesTreeBuilder
