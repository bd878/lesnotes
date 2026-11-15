import type {Message} from '../api/models';
import type {Builder} from './builder';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder';

let messagesTreeTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/messages_tree.mustache')), { encoding: 'utf-8' });
let messagesTreeTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/messages_tree.mustache')), { encoding: 'utf-8' });

class MessagesTreeBuilder extends AbstractPublicBuilder {
	threadID = 0
	list = undefined
	threadPath = undefined

	scripts = ["/public/components/messagesTree/messagesTreeScript.js"]

	addThreadID(id: number) {
		this.threadID = id
		return this
	}

	addList(list: Builder) {
		this.list = list.build()
		return this
	}

	addThreadPath(threadPath: Builder) {
		this.threadPath = threadPath.build()
		return this
	}

	build() {
		return mustache.render(this.isMobile ? messagesTreeTemplateMobile : messagesTreeTemplate, {
			threadID: this.threadID,
		}, {
			list:       this.list,
			threadPath: this.threadPath,
		})
	}
}

export default MessagesTreeBuilder
