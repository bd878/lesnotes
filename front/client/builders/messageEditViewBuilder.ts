import type { Message } from '../api/models';
import type { FileWithMime } from '../types';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let filesListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/files_list.mustache')), { encoding: 'utf-8' });
let filesListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/files_list.mustache')), { encoding: 'utf-8' });

let messageEditFormTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/message_edit_form.mustache')), { encoding: 'utf-8' });
let messageEditFormTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/message_edit_form.mustache')), { encoding: 'utf-8' });

class MessageEditViewBuilder extends AbstractBuilder {
	filesList   = undefined
	message     = undefined

	scripts = [
		"/public/pages/messageEdit/messageEditScript.js"
	]

	addFilesList(files: FileWithMime[]) {
		this.filesList = mustache.render(this.isMobile ? filesListTemplateMobile : filesListTemplate, {
			files: files,
		}, {})
		return this
	}

	addMessage(message: Message) {
		this.message = message
		return this
	}

	build() {
		const search = this.search
		return mustache.render(this.isMobile ? messageEditFormTemplateMobile : messageEditFormTemplate, {
			hasFiles:         this.filesList != undefined,
			message:          this.message,
			filesSummary:     this.i18n("filesSummary"),
			cancelEditHref:   function() { return `/messages/${this.ID}` + search},
			namePlaceholder:  this.i18n("namePlaceholder"),
			titlePlaceholder: this.i18n("titlePlaceholder"),
			textPlaceholder:  this.i18n("textPlaceholder"),
			updateButton:     this.i18n("updateButton"),
			cancelButton:     this.i18n("cancelButton"),
			updateAction:     "/message/update" + this.search,
			domain:           Config.get("domain"),
		}, {
			filesList:        this.filesList,
		})
	}
}

export default MessageEditViewBuilder
