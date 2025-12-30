import type { Thread } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class ThreadEditBuilder extends HomeBuilder {
	async addThreadEditForm(thread: Thread) {
		const template = await readFile(resolve(join(Config.get('basedir'), 
			this.isMobile ? 'templates/home/mobile/thread_edit_form.mustache' : 'templates/home/desktop/thread_edit_form.mustache'
		)), { encoding: 'utf-8' });

		this.threadEditForm = mustache.render(template, {
			ID:               thread.ID,
			private:          thread.private,
			name:             thread.name,
			description:      thread.description,
			cancelEditHref:   `/threads/${thread.ID}` + this.search,
			namePlaceholder:  this.i18n("namePlaceholder"),
			descriptionPlaceholder:  this.i18n("textPlaceholder"),
			updateButton:     this.i18n("updateButton"),
			cancelButton:     this.i18n("cancelButton"),
			updateAction:     "/t/update" + this.search,
			domain:           Config.get("domain"),
		})
	}
}

export default ThreadEditBuilder
