import type { Thread } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder';

let threadViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/thread_view.mustache')), { encoding: 'utf-8' });
let threadViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/thread_view.mustache')), { encoding: 'utf-8' });

class ThreadViewBuilder extends AbstractPublicBuilder {
	thread = undefined

	addThread(thread: Thread) {
		this.thread = thread
		return this
	}

	build() {
		const search = this.search

		return mustache.render(this.isMobile ? threadViewTemplateMobile : threadViewTemplate, {
			thread:           this.thread,
			isAuthed:         this.isAuthed,
			editHref:         function() { return `/editor/threads/${this.ID}` + search; },
			publishAction:    "/t/publish" + search,
			privateAction:    "/t/private" + search,
			domain:           Config.get("domain"),
		})
	}
}

export default ThreadViewBuilder
