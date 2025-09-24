import type { Message } from '../../api/models';
import type { Thread } from '../../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import * as is from '../../third_party/is';
import i18n from '../../i18n';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function home(ctx) {
	let me;
	if (is.empty(ctx.state.me)) {
		console.error("no me")
		ctx.status = 500
		return
	} else {
		me = ctx.state.me.user
	}

	let stack;
	if (is.notEmpty(ctx.state.stack)) {
		if (ctx.state.stack.error.error) {
			console.error(ctx.state.stack.error)
			ctx.body = "error"
			ctx.status = 400;
			return;
		}

		stack = ctx.state.stack.stack
	} else {
		console.error("stack is empty")
		ctx.status = 500
		return
	}

	let message;
	if (is.notEmpty(ctx.state.message)) {
		if (ctx.state.message.error.error) {
			console.error(ctx.state.message.error)
			ctx.body = "error"
			ctx.status = 400;
			return;
		}

		message = ctx.state.message.message
	}

	ctx.set({ "Cache-Control": "no-cache,max-age=0" })

	const builder = new Builder(ctx.userAgent.isMobile)

	if (is.notEmpty(message))
		if (ctx.query.edit)
			await builder.addMessageEditForm(undefined, me.ID, message)
		else
			await builder.addMessageView(undefined, me.ID, message)
	else
		await builder.addNewMessageForm()

	await builder.addMessagesList(undefined, stack)
	await builder.addFilesList(message, ctx.query.edit)
	await builder.addFilesForm()
	await builder.addSearchPath()
	await builder.addSidebar()

	ctx.body = await builder.build(message, ctx.query.edit)
	ctx.status = 200;

	return;
}

class Builder {
	isMobile: boolean = false;
	constructor(isMobile: boolean) {
		this.isMobile = isMobile
	}

	messagesList = undefined;
	async addMessagesList(error: string | undefined, stack: Thread[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/messages_list.mustache' : 'templates/home/desktop/messages_list.mustache'
		)), { encoding: 'utf-8' });

		this.messagesList = mustache.render(template, {
			stack:            stack,
			limit:            14,
			isSingle:         () => stack.length == 1,
			newMessageText:   i18n("newMessageText"),
			noMessagesText:   i18n("noMessagesText"),
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
			namePlaceholder:  i18n("namePlaceholder"),
			titlePlaceholder: i18n("titlePlaceholder"),
			textPlaceholder:  i18n("textPlaceholder"),
			updateButton:     i18n("updateButton"),
			cancelButton:     i18n("cancelButton"),
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
			cancelButton:     i18n("cancelButton"),
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
			titlePlaceholder: i18n("titlePlaceholder"),
			textPlaceholder:  i18n("textPlaceholder"),
			sendButton:       i18n("sendButton"),
		})
	}

	filesList = undefined;
	async addFilesList(message?: Message, editMessage?: boolean) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/files_list.mustache' : 'templates/home/desktop/files_list.mustache'
		)), { encoding: 'utf-8' });

		const options = {
			filesPlaceholder:   i18n("filesPlaceholder"),
			noFiles:            i18n("noFiles"),
			editMessage:        editMessage,
			files:              undefined,
		}

		if (is.notEmpty(message))
			options.files = message.files

		this.filesList = mustache.render(template, options)
	}

	filesForm = undefined;
	async addFilesForm() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/files_form.mustache' : 'templates/home/desktop/files_form.mustache'
		)), { encoding: 'utf-8' });

		this.filesForm = mustache.render(template, {
			filesPlaceholder:    i18n("filesPlaceholder"),
			selectFiles:         i18n("selectFiles"),
		})
	}

	searchPath = undefined;
	async addSearchPath() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/search.mustache' : 'templates/home/desktop/search.mustache'
		)), { encoding: 'utf-8' });

		this.searchPath = mustache.render(template, {
			searchPlaceholder:   i18n("searchPlaceholder"),
		})
	}

	homeSidebar = undefined;
	async addSidebar() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/home_sidebar.mustache' : 'templates/home/desktop/home_sidebar.mustache'
		)), { encoding: 'utf-8' });

		this.homeSidebar = mustache.render(template, {
			logout:           i18n("logout"),
		})
	}

	async build(message?: Message, editMessage?: boolean) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const home = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/home.mustache' : 'templates/home/desktop/home.mustache'
		)), { encoding: 'utf-8' });

		return mustache.render(layout, {
			scripts:  ["/public/pages/home/homeScript.js"],
			manifest: "/public/manifest.json",
			styles:   styles,
			isMobile: this.isMobile ? "true" : "",
		}, {
			content: mustache.render(home, {
				message:     message,
				editMessage: editMessage,
			}, {
				messageEditForm: this.messageEditForm,
				messageView:     this.messageView,
				newMessageForm:  this.newMessageForm,
				messagesList:    this.messagesList,
				homeSidebar:     this.homeSidebar,
				filesForm:       this.filesForm,
				filesList:       this.filesList,
				searchPath:      this.searchPath,
			}),
		});
	}
}

export default home;
