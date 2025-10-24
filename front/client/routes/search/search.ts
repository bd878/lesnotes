import type { Message } from '../../api/models';
import type { File } from '../../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import * as is from '../../third_party/is';
import i18n from '../../i18n';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import Builder from '../builder'

async function search(ctx) {
	let me;
	if (is.empty(ctx.state.me)) {
		console.error("no me")
		ctx.status = 500
		return
	} else {
		me = ctx.state.me
	}

	let messages;
	if (is.notEmpty(ctx.state.search)) {
		if (ctx.state.search.error.error) {
			console.error(ctx.state.search.error)
			ctx.body = "error"
			ctx.status = 400
			return;
		}

		messages = ctx.state.search
	} else {
		console.error("search is empty")
		ctx.status = 500
		return
	}

	ctx.set({ "Cache-Control": "no-cache,max-age=0" })

	const builder = new SearchBuilder(ctx.userAgent.isMobile, ctx.state.lang)

	await builder.addSettings(undefined, ctx.state.lang, me.theme, me.fontSize)
	await builder.addMessagesList(undefined, messages)
	await builder.addFilesList(undefined, undefined)
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	return
}

class SearchBuilder extends Builder {
	messagesList = undefined;
	async addMessagesList(error: string | undefined, list: Message[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/search/mobile/messages_list.mustache' : 'templates/search/desktop/messages_list.mustache'
		)), { encoding: 'utf-8' });

		this.messagesList = mustache.render(template, {
			list:             list,
			isSingle:         () => list.length == 1,
			noMessagesText:   this.i18n("noMessagesText"),
		})
	}

	filesList = undefined;
	async addFilesList(error: string | undefined, list?: File[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/search/mobile/files_list.mustache' : 'templates/search/desktop/files_list.mustache'
		)), { encoding: 'utf-8' });

		const options = {
			filesPlaceholder:   this.i18n("filesPlaceholder"),
			noFiles:            this.i18n("noFiles"),
			files:              undefined,
		}

		if (is.notEmpty(list))
			options.files = list

		this.filesList = mustache.render(template, options)
	}

	async build() {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const search = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/search/mobile/search.mustache' : 'templates/search/desktop/search.mustache'
		)), { encoding: 'utf-8' });

		return mustache.render(layout, {
			scripts:  ["/public/pages/search/searchScript.js"],
			manifest: "/public/manifest.json",
			styles:   styles,
			lang:     this.lang,
			isMobile: this.isMobile ? "true" : "",
		}, {
			footer: this.footer,
			content: mustache.render(search, {}, {
				settings:        this.settings,
				messagesList:    this.messagesList,
				filesList:       this.filesList,
				sidebar:         this.sidebar,
			}),
		});
	}
}

export default search
