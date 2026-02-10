import type { Message } from '../api/models';
import type { SelectedFile } from '../types';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class MessageEditViewBuilder extends HomeBuilder {
	async addMessageEditForm(message?: Message) {
		if (is.empty(message)) {
			return
		}

		const template = await readFile(resolve(join(Config.get('basedir'), 
			this.isMobile ? 'templates/home/mobile/message_edit_form.mustache' : 'templates/home/desktop/message_edit_form.mustache'
		)), { encoding: 'utf-8' });

		this.messageEditForm = mustache.render(template, {
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

	async addFilesSelector(files: SelectedFile[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/files_selector.mustache' : 'templates/home/desktop/files_selector.mustache'
		)), { encoding: 'utf-8' });

		this.filesSelector = mustache.render(template, {
			files:             files,
			defaultFile:       this.i18n("defaultFile"),
		})
	}
}

export default MessageEditViewBuilder
