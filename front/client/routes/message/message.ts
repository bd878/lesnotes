import type { Message } from '../../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import * as is from '../../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import Builder from '../builder';

async function readPublicMessage(ctx) {
	const { message, lang, theme, fontSize } = ctx.state

	ctx.set({ "Cache-Control": "no-cache,max-age=0" })

	if (is.empty(message)) {
		ctx.body = "no message"
		ctx.status = 400

		return
	}

	const builder = new MessageBuilder(ctx.userAgent.isMobile, lang)

	await builder.addSettings(undefined, undefined /* no lang selector */, theme, fontSize)
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build(message, theme, fontSize)
	ctx.status = 200;

	return
}

class MessageBuilder extends Builder {
	sidebar = undefined;
	async addSidebar(query?: string) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_vertical/mobile/sidebar_vertical.mustache' : 'templates/sidebar_vertical/desktop/sidebar_vertical.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template, {
			settingsHeader: this.i18n("settingsHeader")
		}, {
			settings:       this.settings,
		})
	}

	async build(message?: Message, theme?: string, fontSize?: string) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const content = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/message/mobile/message.mustache' : 'templates/message/desktop/message.mustache'
		)), { encoding: 'utf-8' });

		return mustache.render(layout, {
			html:     () => (text, render) => {
				let html = "<html"

				if (theme) html += ` class="${theme}"`;
				if (this.lang) html += ` lang="${this.lang}"`;
				if (fontSize) html += ` data-size="${fontSize}"`
				html += ">"

				return html + render(text) + "</html>"
			},
			manifest:  "/public/manifest.json",
			styles:    styles,
			lang:      this.lang,
			theme:     theme,
			isMobile:  this.isMobile ? "true" : "",
		}, {
			footer:    this.footer,
			content:   mustache.render(content, {
				title:  message.title,
				text:   message.text,
				files:  message.files,
			}, {
				settings:  this.settings,
				sidebar:   this.sidebar,
			})
		})
	}
}

export default readPublicMessage;
