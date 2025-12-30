import type { Thread } from '../../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import * as is from '../../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from '../home/builder';

async function threadView(ctx) {
	console.log("--> threadView")

	const builder = new ThreadViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	if (is.notEmpty(ctx.state.thread)) {
		await builder.addThreadView(ctx.state.thread)
	}
	await builder.addSettings()
	await builder.addMessagesStack(ctx.state.stack)
	await builder.addSearch()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build(undefined, ctx.state.thread)
	ctx.status = 200

	console.log("<-- threadView")
}

class ThreadViewBuilder extends HomeBuilder {
	threadView = undefined;
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

export default threadView;
