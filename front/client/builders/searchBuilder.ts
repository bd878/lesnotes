import type { Message } from '../api/models';
import type { File } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import i18n from '../i18n';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

class SearchBuilder extends AbstractBuilder {
	messagesList = undefined;
	async addMessagesList(list: Message[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/search/mobile/messages_list.mustache' : 'templates/search/desktop/messages_list.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.messagesList = mustache.render(template, {
			list:             list,
			isEmpty:          () => list.length == 0,
			isSingle:         () => list.length == 1,
			emptyListText:    this.i18n("emptyListText"),
			messageHref:      function() { const params = new URLSearchParams(search); params.delete("query"); return `/messages/${this.ID}?` + params.toString(); },
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

	logout = undefined;
	async addLogout() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_vertical/mobile/logout.mustache' : 'templates/sidebar_vertical/desktop/logout.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.logout = mustache.render(template, {
			logout:           this.i18n("logout"),
			logoutHref:       function() { const params = new URLSearchParams(search); params.delete("cwd"); params.delete("id"); /* TODO: delete pagination */ return "/logout?" + params.toString() },
		})
	}

	searchForm = undefined;
	async addSearch() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/search/mobile/search_form.mustache' : 'templates/search/desktop/search_form.mustache'
		)), { encoding: 'utf-8' });

		this.searchForm = mustache.render(template, {
			action:              "/search" + this.search,
			searchPlaceholder:   this.i18n("searchPlaceholder"),
			searchMessages:      this.i18n("search"),
		})
	}

	sidebar = undefined;
	async addSidebar() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_vertical/mobile/sidebar_vertical.mustache' : 'templates/sidebar_vertical/desktop/sidebar_vertical.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.sidebar = mustache.render(template, {
			mainHref:         function() { const params = new URLSearchParams(search); params.delete("query"); return "/home?" + params.toString(); },
			logoutHref:       function() { const params = new URLSearchParams(search); params.delete("query"); return "/logout?" + params.toString(); },
			logout:           this.i18n("logout"),
			settingsHeader:   this.i18n("settingsHeader"),
		}, {
			settings:         this.settings,
			searchForm:       this.searchForm,
			logout:           this.logout,
		})
	}

	async build() {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/layout/mobile/layout.mustache' : 'templates/layout/desktop/layout.mustache'
		)), { encoding: 'utf-8' });
		const search = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/search/mobile/search.mustache' : 'templates/search/desktop/search.mustache'
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

export default SearchBuilder
