import type { Thread } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class ThreadViewBuilder extends HomeBuilder {
	async addThreadView(thread: Thread) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/thread_view.mustache' : 'templates/home/desktop/thread_view.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.threadView = mustache.render(template, {
			ID:               thread.ID,
			description:      thread.description,
			name:             thread.name,
			private:          thread.private,
			newNoteHref:      function() { return "/home" + search; },
			editHref:         function() { return `/editor/threads/${thread.ID}` + search; },
			publishAction:    "/t/publish" + search,
			privateAction:    "/t/private" + search,
			newNoteButton:    this.i18n("newNote"),
			domain:           Config.get("domain"),
		})
	}
}

export default ThreadViewBuilder
