import type { File, Message, Thread, ThreadMessages } from '../api/models';
import type { FileWithMime } from '../types';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import { unwrapPaging } from '../api/models/paging';
import * as is from '../third_party/is';
import i18n from '../i18n';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'

class HomeBuilder extends AbstractBuilder {
	messageEditForm = undefined;
	messageView     = undefined;
	newMessageForm  = undefined;
	threadView      = undefined;
	threadEditForm  = undefined;
	pagination      = undefined;
	filesSelector   = undefined;
	filesForm       = undefined;
	filesView       = undefined;
	filesList       = undefined;
	messagesStack   = undefined;
	searchForm      = undefined;
	logout          = undefined;
	sidebar         = undefined;
	goBack          = undefined;
	controlPanel    = undefined;
	navigation      = undefined;

	async addMessagesStack(stack: ThreadMessages[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/messages_stack.mustache' : 'templates/home/desktop/messages_stack.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search
		const path = this.path

		const close = ((new URLSearchParams(search)).get("close") || "").split(",").map(parseFloat).filter(v => !isNaN(v))

		const limit = parseInt(LIMIT)

		this.messagesStack = mustache.render(template, {
			stack:            stack.map(unwrapPaging),
			limit:            LIMIT,
			isSingle:         () => stack.length == 1,
			openHref:         function() { const params = new URLSearchParams(search); params.set("close", this.closeIDs); return path + "?" + params.toString() },
			closeHref:        function() { const params = new URLSearchParams(search); params.set("close", this.openIDs); return path + "?" + params.toString() },
			newMessageText:   this.i18n("newMessageText"),
			noMessagesText:   this.i18n("noMessagesText"),
			messageHref:      function() { return `/messages/${this.ID}` + search; },
			messageThreadHref: function() { const params = new URLSearchParams(search); params.set("cwd", `${this.ID}`); return path + "?" + params.toString(); },
			viewThreadHref:   function() { return `/threads/${this}` + search; /*context is ID, not thread*/ },
			closeThreadHref:  function() {
				const params = new URLSearchParams(search);
				const set = new Set(close);
				set.add(this.ID);
				(this.parentID == 0 ? params.delete("cwd") : params.set("cwd", `${this.parentID}`));
				params.set("close", Array.from(set).join(","));
				return path + "?" + params.toString();
			},
			rootThreadHref:   function() { const params = new URLSearchParams(search); params.delete("cwd"); return path + "?" + params.toString(); },
			prevPageHref:     function() { const params = new URLSearchParams(search); params.set(this.message.ID || 0, `${limit + this.offset}`); return path + "?" + params.toString(); },
			nextPageHref:     function() { const params = new URLSearchParams(search); params.set(this.message.ID || 0, `${Math.max(0, this.offset - limit)}`); return path + "?" + params.toString(); },
		})
	}

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

	async addControlPanel() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/control_panel.mustache' : 'templates/home/desktop/control_panel.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search;

		this.controlPanel = mustache.render(template, {
			newNoteHref:      function() { return "/home" + search; },
			newNoteButton:    this.i18n("newNote"),
			newFileHref:      function() { return "/files" + search; },
			newFileButton:    this.i18n("newFile"),
		})
	}

	async addNavigation() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/navigation.mustache' : 'templates/home/desktop/navigation.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search;

		this.navigation = mustache.render(template, {
			messagesHref:          function() { return "/home" + search; },
			messagesSection:       this.i18n("messagesSection"),
			filesHref:             function() { return "/files" + search; },
			filesSection:          this.i18n("filesSection"),
		})
	}

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

	async addSidebar() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_vertical/mobile/sidebar_vertical.mustache' : 'templates/sidebar_vertical/desktop/sidebar_vertical.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template, {
			mainHref:         "/home" + this.search,
			settingsHeader:   this.i18n("settingsHeader"),
		}, {
			settings:         this.settings,
			searchForm:       this.searchForm,
			logout:           this.logout,
		})
	}

	async addFilesSelector(files: File[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/files_selector.mustache' : 'templates/home/desktop/files_selector.mustache'
		)), { encoding: 'utf-8' });

		this.filesSelector = mustache.render(template, {
			files:             files,
			defaultFile:       this.i18n("defaultFile"),
		})
	}

	async addFilesView(files: FileWithMime[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/files_view.mustache' : 'templates/home/desktop/files_view.mustache'
		)), { encoding: 'utf-8' });

		this.filesView = mustache.render(template, {
			files:    files,
			imgSrc:   function() { return `/files/v1/read/${this.name}` },
			fileHref: function() { return `/files/v1/download?id=${this.ID}` },
		})
	}

	async build() {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/layout/mobile/layout.mustache' : 'templates/layout/desktop/layout.mustache'
		)), { encoding: 'utf-8' });
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
		}, {
			footer: this.footer,
			content: mustache.render(home, {}, {
				settings:        this.settings,
				messageEditForm: this.messageEditForm,
				messageView:     this.messageView,
				threadView:      this.threadView,
				threadEditForm:  this.threadEditForm,
				newMessageForm:  this.newMessageForm,
				messagesStack:   this.messagesStack,
				sidebar:         this.sidebar,
				pagination:      this.pagination,
				filesList:       this.filesList,
				filesForm:       this.filesForm,
				goBack:          this.goBack,
				filesView:       this.filesView,
				controlPanel:    this.controlPanel,
				navigation:      this.navigation,
			}),
		});
	}
}

export default HomeBuilder
