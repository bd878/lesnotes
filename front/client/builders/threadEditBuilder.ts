import type { Thread } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

let threadEditFormTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/thread_edit_form.mustache')), { encoding: 'utf-8' });
let threadEditFormTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/thread_edit_form.mustache')), { encoding: 'utf-8' });

class ThreadEditBuilder extends HomeBuilder {
	addThreadEditForm(thread: Thread) {
		this.threadEditForm = mustache.render(this.isMobile ? threadEditFormTemplateMobile : threadEditFormTemplate, {
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
