import type { Message } from '../api/models';
import type { SelectedFile } from '../types';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

let messageEditFormTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/message_edit_form.mustache')), { encoding: 'utf-8' });
let messageEditFormTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/message_edit_form.mustache')), { encoding: 'utf-8' });

let filesSelectorTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/files_selector.mustache')), { encoding: 'utf-8' });
let filesSelectorTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/files_selector.mustache')), { encoding: 'utf-8' });

class MessageEditViewBuilder extends HomeBuilder {
	scripts = [
		"/public/pages/home/homeScript.js",
		"/public/pages/messageEdit/messageEditScript.js"
	]

	addMessageEditForm(message?: Message) {
		if (is.empty(message)) {
			return
		}

		this.messageEditForm = mustache.render(this.isMobile ? messageEditFormTemplateMobile : messageEditFormTemplate, {
			ID:               message.ID,
			private:          message.private,
			name:             message.name,
			title:            message.title,
			text:             message.text,
			filesSummary:     this.i18n("filesSummary"),
			cancelEditHref:   `/messages/${message.ID}` + this.search,
			namePlaceholder:  this.i18n("namePlaceholder"),
			titlePlaceholder: this.i18n("titlePlaceholder"),
			textPlaceholder:  this.i18n("textPlaceholder"),
			updateButton:     this.i18n("updateButton"),
			cancelButton:     this.i18n("cancelButton"),
			updateAction:     "/m/update" + this.search,
			domain:           Config.get("domain"),
		}, {
			filesSelector:    this.filesSelector,
		})
	}

	addFilesSelector(files: SelectedFile[]) {
		this.filesSelector = mustache.render(this.isMobile ? filesSelectorTemplateMobile : filesSelectorTemplate, {
			files:             files,
			defaultFile:       this.i18n("defaultFile"),
		})
	}
}

export default MessageEditViewBuilder
