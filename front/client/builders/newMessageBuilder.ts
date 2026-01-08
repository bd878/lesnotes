import Config from 'config';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class NewMessageBuilder extends HomeBuilder {
	async addNewMessageForm(thread?: number) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/new_message_form.mustache' : 'templates/home/desktop/new_message_form.mustache'
		)), { encoding: 'utf-8' });

		this.newMessageForm = mustache.render(template, {
			titlePlaceholder: this.i18n("titlePlaceholder"),
			textPlaceholder:  this.i18n("textPlaceholder"),
			sendButton:       this.i18n("sendButton"),
			sendAction:       "/send" + this.search,
			thread:           thread || 0,
		}, {
			filesInput: this.filesInput,
		})
	}
}

export default NewMessageBuilder
