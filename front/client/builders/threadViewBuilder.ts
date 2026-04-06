import type { Thread } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let threadViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/thread_view.mustache')), { encoding: 'utf-8' });
let threadViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/thread_view.mustache')), { encoding: 'utf-8' });

class ThreadViewBuilder extends AbstractBuilder {
	thread = undefined

	addThread(thread: Thread) {
		this.thread = thread
	}

	build() {
		const search = this.search

		return mustache.render(this.isMobile ? threadViewTemplateMobile : threadViewTemplate, {
			thread:           this.thread,
			newNoteHref:      function() { return "/home" + search; },
			editHref:         function() { return `/editor/threads/${this.ID}` + search; },
			publishAction:    "/t/publish" + search,
			privateAction:    "/t/private" + search,
			newNoteButton:    this.i18n("newNote"),
			domain:           Config.get("domain"),
		})
	}
}

export default ThreadViewBuilder
