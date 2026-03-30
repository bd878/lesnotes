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

class MessageHeaderBuilder extends AbstractBuilder {
	messagePath = undefined

	addMessagePath(path: string) {
		this.messagePath = mustache.render(this.isMobile ? messagePathTemplateMobile : messagePathTemplate, {
			path: path,
		})
	}

	build() {
		return mustache.render(this.isMobile ? messageHeaderTemplateMobile : messageHeaderTemplate, {}, {
			messagePath: this.messagePath,
		})
	}
}

export default MessageHeaderBuilder
