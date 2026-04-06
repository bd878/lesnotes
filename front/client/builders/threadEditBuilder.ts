import type { Thread } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let threadEditFormTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/thread_edit_form.mustache')), { encoding: 'utf-8' });
let threadEditFormTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/thread_edit_form.mustache')), { encoding: 'utf-8' });

class ThreadEditBuilder extends AbstractBuilder {
	thread = undefined

	addThread(thread: Thread) {
		this.thread = thread
		return this
	}

	build() {
		const search = this.search

		return mustache.render(this.isMobile ? threadEditFormTemplateMobile : threadEditFormTemplate, {
			thread:           this.thread,
			cancelEditHref:   function() { return `/threads/${this.ID}` + search },
			namePlaceholder:  this.i18n("namePlaceholder"),
			descriptionPlaceholder:  this.i18n("descriptionPlaceholder"),
			titlePlaceholder: this.i18n("titlePlaceholder"),
			updateButton:     this.i18n("updateButton"),
			cancelButton:     this.i18n("cancelButton"),
			updateAction:     "/t/update" + this.search,
			domain:           Config.get("domain"),
		})
	}
}

export default ThreadEditBuilder
