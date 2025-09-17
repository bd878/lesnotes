import type { Message } from '../api/models';
import type { Thread } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import i18n from '../i18n';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

/**
 * Renders home page
 * 
 * Example: /home?thread=123&id=111&limit=5&offset=10
 * thread - message thread id
 * id - message to render
 * limit - limit messages of final thread to load
 * offset - messages offset of final thread
 * 
 * @param {[type]} ctx context
 */
async function home(ctx) {
	const token = ctx.cookies.get("token")

	ctx.set({ "Cache-Control": "no-cache,max-age=0" })

	const resp = await api.authJson(token)
	console.log(`[home]: auth response`, JSON.stringify(resp))
	if (resp.error.error || resp.expired) {
		ctx.redirect("/login")
		ctx.status = 302
		return
	}

	const me = await api.getMeJson(token)
	if (me.error.error) {
		ctx.redirect("/login")
		ctx.status = 302
		return
	}

	const edit = ctx.query.edit
	const limit = parseInt(ctx.query.limit) || 10
	const offset = parseInt(ctx.query.offset) || 0
	const threadID = parseInt(ctx.query.thread) || 0
	const id = parseInt(ctx.query.id) || 0

	const stack = await api.readStackJson(token, threadID, id, 10)
	if (stack.error.error) {
		console.log(stack.error)
		ctx.body = await renderError("failed to load messages stack");
		ctx.status = 400;
		return;
	}

	let message;
	if (id != 0) {
		message = await api.readMessageJson(token, 0, id)
		if (message.error.error) {
			ctx.body = await renderError("failed to load message")
			ctx.status = 400;
			return;
		}

		ctx.body = await renderBody(stack.stack, me.user.ID, message.message, edit)
		ctx.status = 200;

		return
	}

	ctx.body = await renderBody(stack.stack, me.user.ID, undefined, edit)
	ctx.status = 200;

	return;
}

async function renderError(err: string): Promise<string> {
	const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles.css')), { encoding: 'utf-8' });
	const home = await readFile(resolve(join(Config.get('basedir'), 'templates/home.mustache')), { encoding: 'utf-8' });
	const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
	const messageEditForm = await readFile(resolve(join(Config.get('basedir'), 'templates/message_edit_form.mustache')), { encoding: 'utf-8' });
	const messageView = await readFile(resolve(join(Config.get('basedir'), 'templates/message_view.mustache')), { encoding: 'utf-8' });
	const newMessageForm = await readFile(resolve(join(Config.get('basedir'), 'templates/new_message_form.mustache')), { encoding: 'utf-8' });
	const homeSidebar = await readFile(resolve(join(Config.get('basedir'), 'templates/home_sidebar.mustache')), { encoding: 'utf-8' });
	const messagesList = await readFile(resolve(join(Config.get('basedir'), 'templates/messages_list.mustache')), { encoding: 'utf-8' });

	const content = mustache.render(home, {
		error:    err,
		domain:   Config.get("domain"),
		send:     i18n("send"),
		logout:   i18n("logout"),
		filesPlaceholder: i18n("filesPlaceholder"),
		newMessageText: i18n("newMessageText"),
		selectFiles: i18n("selectFiles"),
		search:   i18n("search"),
		delete:   i18n("delete"),
		edit:     i18n("edit"),
		publish:  i18n("publish"),
		privateText:  i18n("private"),
		update:   i18n("update"),
		cancel:        i18n("cancel"),
		noFiles:        i18n("noFiles"),
		namePlaceholder:  i18n("namePlaceholder"),
		titlePlaceholder: i18n("titlePlaceholder"),
		textPlaceholder:  i18n("textPlaceholder"),
	}, {
		messageEditForm,
		messageView,
		newMessageForm,
		homeSidebar,
		messagesList,
	})

	return mustache.render(layout, {
		scripts:  ["/public/homeScript.js"],
		manifest: "/public/manifest.json",
		styles:   styles,
	}, {
		content,
	});
}

async function renderBody(stack: Thread[], userID: number, message?: Message, editMessage?: boolean): Promise<string> {
	const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles.css')), { encoding: 'utf-8' });
	const home = await readFile(resolve(join(Config.get('basedir'), 'templates/home.mustache')), { encoding: 'utf-8' });
	const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
	const messageEditForm = await readFile(resolve(join(Config.get('basedir'), 'templates/message_edit_form.mustache')), { encoding: 'utf-8' });
	const messageView = await readFile(resolve(join(Config.get('basedir'), 'templates/message_view.mustache')), { encoding: 'utf-8' });
	const newMessageForm = await readFile(resolve(join(Config.get('basedir'), 'templates/new_message_form.mustache')), { encoding: 'utf-8' });
	const homeSidebar = await readFile(resolve(join(Config.get('basedir'), 'templates/home_sidebar.mustache')), { encoding: 'utf-8' });
	const messagesList = await readFile(resolve(join(Config.get('basedir'), 'templates/messages_list.mustache')), { encoding: 'utf-8' });

	const content = mustache.render(home, {
		stack:    stack,
		message:  message,
		userID:   userID,
		domain:   Config.get("domain"),
		send:     i18n("send"),
		logout:   i18n("logout"),
		filesPlaceholder: i18n("filesPlaceholder"),
		newMessageText: i18n("newMessageText"),
		selectFiles: i18n("selectFiles"),
		search:   i18n("search"),
		delete:   i18n("delete"),
		edit:     i18n("edit"),
		editMessage: editMessage,
		publish:  i18n("publish"),
		privateText:  i18n("private"),
		update:   i18n("update"),
		cancel:        i18n("cancel"),
		noFiles:        i18n("noFiles"),
		namePlaceholder:  i18n("namePlaceholder"),
		titlePlaceholder: i18n("titlePlaceholder"),
		textPlaceholder:  i18n("textPlaceholder"),
	}, {
		messageEditForm,
		messageView,
		newMessageForm,
		messagesList,
		homeSidebar,
	})

	return mustache.render(layout, {
		scripts:  ["/public/homeScript.js"],
		manifest: "/public/manifest.json",
		styles:   styles,
	}, {
		content,
	});
}

export default home;
