import type { Thread, ThreadMessages } from '../../api/models';
import { unwrapPaging } from '../../api/models/paging';
import Config from 'config';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import Builder from '../builder';

class PublicThreadBuilder extends Builder {
	sidebar = undefined;
	async addSidebar() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_vertical/mobile/sidebar_vertical.mustache' : 'templates/sidebar_vertical/desktop/sidebar_vertical.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template, {
			settingsHeader: this.i18n("settingsHeader")
		}, {
			settings:       this.settings,
			searchForm:     this.searchForm,
		})
	}

	messagesList = undefined
	async addMessagesList(name: string, messages: ThreadMessages) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/thread/mobile/messages_list.mustache' : 'templates/thread/desktop/messages_list.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search
		const path = this.path

		const limit = parseInt(LIMIT)

		this.messagesList = mustache.render(template, {
			messages:         unwrapPaging(messages),
			limit:            limit,
			isSingle:         () => messages.messages.length == 1,
			messageHref:      function() { return `/t/${name}/${this.name}` + search; },
			prevPageHref:     function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${limit + this.offset}`); return path + "?" + params.toString(); },
			nextPageHref:     function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${Math.max(0, this.offset - limit)}`); return path + "?" + params.toString(); },
		})
	}

	searchForm = undefined;
	async addSearch() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/search_form/mobile/search_form.mustache' : 'templates/search_form/desktop/search_form.mustache'
		)), { encoding: 'utf-8' });

		this.searchForm = mustache.render(template, {
			action:              "/search" + this.search,
			searchPlaceholder:   this.i18n("searchPlaceholder"),
			searchMessages:      this.i18n("search"),
		})
	}

	async build(thread?: Thread) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const content = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/thread/mobile/thread.mustache' : 'templates/thread/desktop/thread.mustache'
		)), { encoding: 'utf-8' });

		const theme = this.theme
		const fontSize = this.fontSize

		return mustache.render(layout, {
			html: () => (text, render) => {
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
			theme:    this.theme,
			isMobile: this.isMobile ? "true" : "",
		}, {
			footer: this.footer,
			content: mustache.render(content, {
				thread:        thread,
			}, {
				settings:      this.settings,
				sidebar:       this.sidebar,
				messagesList:  this.messagesList,
			})
		})
	}
}

export default PublicThreadBuilder
