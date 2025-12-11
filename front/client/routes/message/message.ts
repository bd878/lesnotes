import type { Message } from '../../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import * as is from '../../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import Builder from '../builder';

async function readPublicMessage(ctx) {
	const { message } = ctx.state

	ctx.set({ "Cache-Control": "no-cache,max-age=0" })

	const builder = new MessageBuilder(ctx.userAgent.isMobile, ctx.state.lang)

	if (is.notEmpty(message))
		await builder.addMessage(message)

	ctx.body = await builder.build(ctx.state.theme, ctx.state.fontSize)
	ctx.status = 200;

	return
}

class MessageBuilder {
	isMobile: boolean = false;
	constructor(isMobile: boolean) {
		this.isMobile = isMobile
	}

	message = undefined;
	async addMessage(message?: Message) {
		if (is.empty(message))
			return;

		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/message/mobile/message.mustache' : 'templates/message/desktop/message.mustache'
		)), { encoding: 'utf-8' });

		this.message = mustache.render(template, {
			title:  message.title,
			text:   message.text,
			files:  message.files,
		})
	}

	async build(theme?: string, fontSize?: string) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });

		return mustache.render(layout, {
			html:     () => (text, render) => {
				let html = "<html"

				if (theme) html += ` class="${theme}"`;
				if (this.lang) html += ` lang="${this.lang}"`;
				if (fontSize) html += ` data-size="${fontSize}"`
				html += ">"

				return html + render(text) + "</html>"
			},
			scripts:   ["/public/pages/message/messageScript.js"],
			manifest:  "/public/manifest.json",
			styles:    styles,
			isMobile:  this.isMobile ? "true" : "",
		}, {
			content:   this.message,
		})
	}
}

export default readPublicMessage;
