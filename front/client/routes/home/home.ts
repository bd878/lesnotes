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

async function home(ctx) {
	const { me, stack, message } = ctx.state

	ctx.set({ "Cache-Control": "no-cache,max-age=0" })

	const builder = new HomeBuilder(ctx.userAgent.isMobile, ctx.state.lang)

	switch (ctx.state.editorMode) {
	case "view":
		await builder.addMessageView(undefined, me.ID, message)
		break;
	case "edit":
		await builder.addMessageEditForm(undefined, me.ID, message)
		break;
	case "new-message":
		await builder.addNewMessageForm()
		break
	default:
		console.error("unknown editor mode")
		ctx.status = 500
		return
	}

	await builder.addSettings(undefined, ctx.state.lang, me.theme, me.fontSize)
	await builder.addMessagesList(undefined, stack)
	await builder.addFilesList(message, ctx.query.edit)
	await builder.addFilesForm(message, ctx.query.edit)
	await builder.addSearch()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build(message, ctx.query.edit, me.theme, me.fontSize)
	ctx.status = 200;

	return;
}

class HomeBuilder extends Builder {
	messagesList = undefined;
	async addMessagesList(error: string | undefined, stack: Thread[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/messages_list.mustache' : 'templates/home/desktop/messages_list.mustache'
		)), { encoding: 'utf-8' });

		// TODO: make messages a link, not a button
		// generate hrefs here
		this.messagesList = mustache.render(template, {
			stack:            stack,
			limit:            LIMIT,
			isSingle:         () => stack.length == 1,
			newMessageText:   this.i18n("newMessageText"),
			noMessagesText:   this.i18n("noMessagesText"),
		})
	}

	messageEditForm = undefined;
	async addMessageEditForm(error: string | undefined, userID: number, message?: Message) {
		if (is.empty(message))
			return

		const template = await readFile(resolve(join(Config.get('basedir'), 
			this.isMobile ? 'templates/home/mobile/message_edit_form.mustache' : 'templates/home/desktop/message_edit_form.mustache'
		)), { encoding: 'utf-8' });

		this.messageEditForm = mustache.render(template, {
			ID:               message.ID,
			private:          message.private,
			name:             message.name,
			title:            message.title,
			text:             message.text,
			namePlaceholder:  this.i18n("namePlaceholder"),
			titlePlaceholder: this.i18n("titlePlaceholder"),
			textPlaceholder:  this.i18n("textPlaceholder"),
			updateButton:     this.i18n("updateButton"),
			cancelButton:     this.i18n("cancelButton"),
			userID:           userID,
			domain:           Config.get("domain"),
		})
	}

	messageView = undefined;
	async addMessageView(error: string | undefined, userID: number, message?: Message) {
		if (is.empty(message))
			return

		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/message_view.mustache' : 'templates/home/desktop/message_view.mustache'
		)), { encoding: 'utf-8' });

		this.messageView = mustache.render(template, {
			ID:               message.ID,
			title:            message.title,
			text:             message.text,
			name:             message.name,
			private:          message.private,
			newNoteButton:    this.i18n("newNote"),
			userID:           userID,
			domain:           Config.get("domain"),
		})
	}

	newMessageForm = undefined;
	async addNewMessageForm() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/new_message_form.mustache' : 'templates/home/desktop/new_message_form.mustache'
		)), { encoding: 'utf-8' });

		this.newMessageForm = mustache.render(template, {
			titlePlaceholder: this.i18n("titlePlaceholder"),
			textPlaceholder:  this.i18n("textPlaceholder"),
			sendButton:       this.i18n("sendButton"),
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

	async build(message?: Message, editMessage?: boolean, theme?: string, fontSize?: string) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const home = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/home.mustache' : 'templates/home/desktop/home.mustache'
		)), { encoding: 'utf-8' });

		return mustache.render(layout, {
			html:   () => (text, render) => {
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
				editMessage: editMessage,
			}, {
				settings:        this.settings,
				messageEditForm: this.messageEditForm,
				messageView:     this.messageView,
				newMessageForm:  this.newMessageForm,
				messagesList:    this.messagesList,
				sidebar:         this.sidebar,
				filesList:       this.filesList,
				filesForm:       this.filesForm,
			}),
		});
	}
}

export default home;
