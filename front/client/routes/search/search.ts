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
	console.log("--> search")

	let messages;
	if (is.notEmpty(ctx.state.searchPath)) {
		if (ctx.state.search.error.error) {
			console.error(ctx.state.search.error)
			ctx.body = "error"
			ctx.status = 400
			return;
		}

		if (ctx.state.searchPath.error.error) {
			console.error(ctx.state.searchPath.error)
			ctx.body = "error"
			ctx.status = 400
			return;
		}

		messages = ctx.state.searchPath.messages
	} else {
		console.error("search is empty")
		ctx.status = 500
		return
	}

	const builder = new SearchBuilder(ctx.userAgent.isMobile, ctx.state.lang)

	await builder.addSettings(ctx.state.lang, ctx.state.theme, ctx.state.fontSize)
	await builder.addMessagesList(messages)
	await builder.addFilesList()
	await builder.addSearch()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build(ctx.state.theme, ctx.state.fontSize)
	ctx.status = 200

	console.log("<-- search")
}

class SearchBuilder extends Builder {
	messagesList = undefined;
	async addMessagesList(list: Message[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/search/mobile/messages_list.mustache' : 'templates/search/desktop/messages_list.mustache'
		)), { encoding: 'utf-8' });

		this.messagesList = mustache.render(template, {
			list:             list,
			isEmpty:          () => list.length == 0,
			isSingle:         () => list.length == 1,
			emptyListText:    this.i18n("emptyListText"),
		})
	}

	filesList = undefined;
	async addFilesList(list?: File[]) {
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

	searchForm = undefined;
	async addSearch() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/search/mobile/search_form.mustache' : 'templates/search/desktop/search_form.mustache'
		)), { encoding: 'utf-8' });

		this.searchForm = mustache.render(template, {
			searchPlaceholder:   this.i18n("searchPlaceholder"),
			searchMessages:      this.i18n("search"),
		})
	}

	sidebar = undefined;
	async addSidebar() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_vertical/mobile/sidebar_vertical.mustache' : 'templates/sidebar_vertical/desktop/sidebar_vertical.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template, {
			logout:           this.i18n("logout"),
			settingsHeader:   this.i18n("settingsHeader"),
		}, {
			settings:         this.settings,
			searchForm:       this.searchForm,
		})
	}

	async build(theme?: string, fontSize?: string) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const search = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/search/mobile/search.mustache' : 'templates/search/desktop/search.mustache'
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
			manifest: "/public/manifest.json",
			styles:   styles,
			lang:     this.lang,
			isMobile: this.isMobile ? "true" : "",
		}, {
			footer: this.footer,
			content: mustache.render(search, {}, {
				messagesList:    this.messagesList,
				filesList:       this.filesList,
				sidebar:         this.sidebar,
			}),
		});
	}
}

export default search
