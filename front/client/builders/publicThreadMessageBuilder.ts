import type { Message } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import PublicThreadBuilder from './publicThreadBuilder'

class PublicThreadMessageBuilder extends PublicThreadBuilder {
	async addMessageView(message: Message) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/thread/mobile/message_view.mustache' : 'templates/thread/desktop/message_view.mustache'
		)), { encoding: 'utf-8' });

		this.messageView = mustache.render(template, {
			message: message,
		})
	}
}

export default PublicThreadMessageBuilder
