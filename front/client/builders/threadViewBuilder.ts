import type { Thread } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

let threadViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/thread_view.mustache')), { encoding: 'utf-8' });
let threadViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/thread_view.mustache')), { encoding: 'utf-8' });

class ThreadViewBuilder extends HomeBuilder {
	addThreadView(thread: Thread) {
		const search = this.search

		this.threadView = mustache.render(this.isMobile ? threadViewTemplate : threadViewTemplateMobile, {
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
