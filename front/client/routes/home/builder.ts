import type { Message } from '../../api/models';
import type { Thread } from '../../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import * as is from '../../third_party/is';
import i18n from '../../i18n';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import Builder from '../builder'

class HomeBuilder extends Builder {
	messageEditForm = undefined;
	messageView = undefined;
	threadView = undefined;
	threadEditForm = undefined;

	messagesList = undefined;
	async addMessagesList(stack: Thread[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/messages_list.mustache' : 'templates/home/desktop/messages_list.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search
		const path = this.path

		const limit = parseInt(LIMIT)

		this.messagesList = mustache.render(template, {
			stack:            stack,
			limit:            LIMIT,
			isSingle:         () => stack.length == 1,
			newMessageText:   this.i18n("newMessageText"),
			noMessagesText:   this.i18n("noMessagesText"),
			messageHref:      function() { return `/messages/${this.ID}` + search; },
			messageThreadHref: function() { const params = new URLSearchParams(search); params.set("cwd", `${this.ID}`); return path + "?" + params.toString(); },
			viewThreadHref:   function() { return `/threads/${this}` + search; /*context is ID, not thread*/ },
			closeThreadHref:  function() { const params = new URLSearchParams(search); params.set("cwd", `${this}` /*context is parentID*/); return path + "?" + params.toString(); },
			closeRootChildThreadHref: function() { const params = new URLSearchParams(search); params.delete("cwd"); return path + "?" + params.toString(); },
			rootThreadHref:   function() { const params = new URLSearchParams(search); params.delete("cwd"); return path + "?" + params.toString(); },
			prevPageHref:     function() { const params = new URLSearchParams(search); params.set(this.threadID || 0, `${limit + this.offset}`); return path + "?" + params.toString(); },
			nextPageHref:     function() { const params = new URLSearchParams(search); params.set(this.threadID || 0, `${Math.max(0, this.offset - limit)}`); return path + "?" + params.toString(); },
		})
	}

	filesList = undefined;
	async addFilesList(message?: Message, editMessage?: boolean) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/files_list.mustache' : 'templates/home/desktop/files_list.mustache'
		)), { encoding: 'utf-8' });

		const options = {
			noFiles:            this.i18n("noFiles"),
			editMessage:        editMessage,
			files:              undefined,
		}

		if (is.notEmpty(message))
			options.files = message.files

		this.filesList = mustache.render(template, options)
	}

	filesForm = undefined;
	async addFilesForm(message?: Message, editMessage?: boolean) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/files_form.mustache' : 'templates/home/desktop/files_form.mustache'
		)), { encoding: 'utf-8' });

		const options = {
			noFiles:            this.i18n("noFiles"),
			editMessage:        editMessage,
			files:              undefined,
		}

		if (is.notEmpty(message))
			options.files = message.files

		this.filesForm = mustache.render(template, options)
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

	sidebar = undefined;
	async addSidebar(search: string) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_vertical/mobile/sidebar_vertical.mustache' : 'templates/sidebar_vertical/desktop/sidebar_vertical.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template, {
			mainHref:         "/home" + search,
			logout:           this.i18n("logout"),
			logoutHref:       function() { const params = new URLSearchParams(search); params.delete("cwd"); params.delete("id"); /* TODO: delete pagination */ return "/logout?" + params.toString() },
			settingsHeader:   this.i18n("settingsHeader"),
		}, {
			settings:         this.settings,
			searchForm:       this.searchForm,
		})
	}

	newMessageForm = undefined;
	async addNewMessageForm(thread?: number) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/new_message_form.mustache' : 'templates/home/desktop/new_message_form.mustache'
		)), { encoding: 'utf-8' });

		this.newMessageForm = mustache.render(template, {
			titlePlaceholder: this.i18n("titlePlaceholder"),
			textPlaceholder:  this.i18n("textPlaceholder"),
			sendButton:       this.i18n("sendButton"),
			sendAction:       "/send" + this.search,
			thread:           thread || 0,
		})
	}

	async build(message?: Message, thread?: Thread) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const home = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/home.mustache' : 'templates/home/desktop/home.mustache'
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
			scripts:  ["/public/pages/home/homeScript.js"],
			manifest: "/public/manifest.json",
			styles:   styles,
			lang:     this.lang,
			theme:    theme,
			isMobile: this.isMobile ? "true" : "",
		}, {
			footer: this.footer,
			content: mustache.render(home, {
				message:     message,
				thread:      thread,
			}, {
				settings:        this.settings,
				messageEditForm: this.messageEditForm,
				messageView:     this.messageView,
				threadView:      this.threadView,
				threadEditForm:  this.threadEditForm,
				newMessageForm:  this.newMessageForm,
				messagesList:    this.messagesList,
				sidebar:         this.sidebar,
				filesList:       this.filesList,
				filesForm:       this.filesForm,
			}),
		});
	}
}

export default HomeBuilder
