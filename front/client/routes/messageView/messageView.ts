import type { Message } from '../../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import * as is from '../../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from '../home/builder';

async function messageView(ctx) {
	console.log("--> messageView")

	const builder = new MessageViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addMessageView(ctx.state.me.ID, ctx.state.message)
	await builder.addSettings()
	await builder.addMessagesList(ctx.state.stack)
	await builder.addSearch()
	await builder.addSidebar(ctx.search)
	await builder.addFooter()

	ctx.body = await builder.build(ctx.state.message)
	ctx.status = 200

	console.log("<-- messageView")
}

class MessageViewBuilder extends HomeBuilder {
	messageView = undefined;
	async addMessageView(userID: number, message?: Message) {
		if (is.empty(message))
			return

		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/message_view.mustache' : 'templates/home/desktop/message_view.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.messageView = mustache.render(template, {
			ID:               message.ID,
			title:            message.title,
			text:             message.text,
			name:             message.name,
			private:          message.private,
			newNoteHref:      function() { return "/home" + search; },
			editHref:         function() { return `/editor/messages/${message.ID}` + search; },
			deleteAction:     "/m/delete" + search,
			publishAction:    "/m/publish" + search,
			privateAction:    "/m/private" + search,
			newNoteButton:    this.i18n("newNote"),
			userID:           userID,
			domain:           Config.get("domain"),
		})
	}
}

export default messageView;
