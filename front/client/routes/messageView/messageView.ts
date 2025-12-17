import type { Message } from '../../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import * as is from '../../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from '../home/builder';

async function readMessageView(ctx) {
	console.log("--> message view")

	const { me, stack, message, lang, theme, fontSize } = ctx.state

	ctx.set({ "Cache-Control": "no-cache,max-age=0" })

	if (is.empty(message)) {
		ctx.body = "no message"
		ctx.status = 400

		return
	}

	const builder = new MessageViewBuilder(ctx.userAgent.isMobile, lang)

	await builder.addMessageView(undefined, me.ID, message)
	await builder.addSettings(undefined, lang, theme, fontSize)
	await builder.addMessagesList(undefined, stack)
	await builder.addSearch()
	await builder.addSidebar(ctx.search)
	await builder.addFooter()

	ctx.body = await builder.build(message, theme, fontSize)
	ctx.status = 200

	console.log("<-- message view")

	return;
}

class MessageViewBuilder extends HomeBuilder {
	messageView = undefined;
	async addMessageView(error: string | undefined, userID: number, message?: Message) {
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
			newNoteHref:      function() { const p = new URLSearchParams(search); p.delete("id"); return "/home?" + p.toString(); },
			editHref:         function() { const p = new URLSearchParams(search); p.set("edit", "1"); return "/home?" + p.toString(); },
			deleteAction:     "/delete" + search,
			publishAction:    "/publish" + search,
			privateAction:    "/private" + search,
			newNoteButton:    this.i18n("newNote"),
			userID:           userID,
			domain:           Config.get("domain"),
		})
	}
}

export default readMessageView;
