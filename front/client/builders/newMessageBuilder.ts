import Config from 'config';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

let newMessageFormTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/new_message_form.mustache')), { encoding: 'utf-8' });
let newMessageFormTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/new_message_form.mustache')), { encoding: 'utf-8' });

class NewMessageBuilder extends HomeBuilder {
	addNewMessageForm(thread?: number) {
		this.newMessageForm = mustache.render(this.isMobile ? newMessageFormTemplateMobile : newMessageFormTemplate, {
			filesSummary:     this.i18n("filesSummary"),
			titlePlaceholder: this.i18n("titlePlaceholder"),
			textPlaceholder:  this.i18n("textPlaceholder"),
			sendButton:       this.i18n("sendButton"),
			sendAction:       "/send" + this.search,
			thread:           thread || 0,
		}, {})
	}
}

export default NewMessageBuilder
