import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let messagePathTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/message_path.mustache')), { encoding: 'utf-8' });
let messagePathTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/message_path.mustache')), { encoding: 'utf-8' });

let messageHeaderTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/message_header.mustache')), { encoding: 'utf-8' });
let messageHeaderTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/message_header.mustache')), { encoding: 'utf-8' });

let messageThreadSwitcherTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/message_thread_switcher.mustache')), { encoding: 'utf-8' });
let messageThreadSwitcherTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/message_thread_switcher.mustache')), { encoding: 'utf-8' });

class MessageHeaderBuilder extends AbstractBuilder {
	messagePath = undefined
	switcher = undefined

	addMessagePath(path: string) {
		this.messagePath = mustache.render(this.isMobile ? messagePathTemplateMobile : messagePathTemplate, {
			path: path,
		})
	}

	addMessageLink(messageID: number) {
		this.switcher = mustache.render(this.isMobile ? messageThreadSwitcherTemplateMobile : messageThreadSwitcherTemplate, {
			href: `/messages/${messageID}` + this.search,
			title: this.i18n("note"),
		})
	}

	addThreadLink(threadID: number) {
		this.switcher = mustache.render(this.isMobile ? messageThreadSwitcherTemplateMobile : messageThreadSwitcherTemplate, {
			href: `/threads/${threadID}` + this.search,
			title: this.i18n("thread"),
		})
	}

	build() {
		return mustache.render(this.isMobile ? messageHeaderTemplateMobile : messageHeaderTemplate, {}, {
			messagePath: this.messagePath,
			switcher:    this.switcher,
		})
	}
}

export default MessageHeaderBuilder
