import Config from 'config';
import type { FileWithMime } from '../types';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let filesListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/files_list.mustache')), { encoding: 'utf-8' });
let filesListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/files_list.mustache')), { encoding: 'utf-8' });

let newMessageFormTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/new_message_form.mustache')), { encoding: 'utf-8' });
let newMessageFormTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/new_message_form.mustache')), { encoding: 'utf-8' });

class NewMessageBuilder extends AbstractBuilder {
	filesList = undefined
	threadID  = 0

	addFilesList(files: FileWithMime[]) {
		this.filesList = mustache.render(this.isMobile ? filesListTemplateMobile : filesListTemplate, {
			files: files,
		}, {})
		return this
	}

	addThreadID(threadID: number) {
		this.threadID = threadID
		return this
	}

	build() {
		return mustache.render(this.isMobile ? newMessageFormTemplateMobile : newMessageFormTemplate, {
			hasFiles:         this.filesList != undefined,
			filesSummary:     this.i18n("filesSummary"),
			titlePlaceholder: this.i18n("titlePlaceholder"),
			textPlaceholder:  this.i18n("textPlaceholder"),
			sendButton:       this.i18n("sendButton"),
			sendAction:       "/send" + this.search,
			thread:           this.threadID,
		}, {
			filesList:        this.filesList,
		})
	}
}

export default NewMessageBuilder
