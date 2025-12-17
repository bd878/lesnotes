import type { Message } from '../../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import * as is from '../../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from '../home/builder';

async function readMessageEdit(ctx) {
	console.log("--> message edit")

	const { me, stack, message, lang, theme, fontSize } = ctx.state

	ctx.set({ "Cache-Control": "no-cache,max-age=0" })

	if (is.empty(message)) {
		ctx.body = "no message"
		ctx.status = 400

		return
	}

	const builder = new MessageEditViewBuilder(ctx.userAgent.isMobile, lang)

	await builder.addMessageEditForm(undefined, me.ID, message)
	await builder.addSettings(undefined, lang, theme, fontSize)
	await builder.addMessagesList(undefined, stack)
	await builder.addSearch()
	await builder.addSidebar(ctx.search)
	await builder.addFooter()

	ctx.body = await builder.build(message, theme, fontSize)
	ctx.status = 200

	console.log("<-- message edit")

	return;
}

class MessageEditViewBuilder extends HomeBuilder {
	messageEditForm = undefined;
	async addMessageEditForm(error: string | undefined, userID: number, message?: Message) {
		if (is.empty(message))
			return

		const template = await readFile(resolve(join(Config.get('basedir'), 
			this.isMobile ? 'templates/home/mobile/message_edit_form.mustache' : 'templates/home/desktop/message_edit_form.mustache'
		)), { encoding: 'utf-8' });

		const params = new URLSearchParams(this.search)
		params.delete("edit")

		this.messageEditForm = mustache.render(template, {
			ID:               message.ID,
			private:          message.private,
			name:             message.name,
			title:            message.title,
			text:             message.text,
			cancelEditHref:   "/home?" + params.toString(),
			namePlaceholder:  this.i18n("namePlaceholder"),
			titlePlaceholder: this.i18n("titlePlaceholder"),
			textPlaceholder:  this.i18n("textPlaceholder"),
			updateButton:     this.i18n("updateButton"),
			cancelButton:     this.i18n("cancelButton"),
			updateAction:     "/update" + this.search,
			userID:           userID,
			domain:           Config.get("domain"),
		})
	}
}

export default readMessageEdit;
