import type { Message } from '../../api/models';
import type { Thread } from '../../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import i18n from '../../i18n';
import * as is from '../../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function renderMobile(ctx) {
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
			ctx.body = await renderError("failed to load messages stack");
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
			ctx.body = await renderError("failed to load message")
			ctx.status = 400;
			return;
		}

		message = ctx.state.message.message
	}

	ctx.set({ "Cache-Control": "no-cache,max-age=0" })

	ctx.body = await renderBody(stack, me.ID, message, ctx.query.edit)
	ctx.status = 200;

	return;
}

async function renderError(err: string): Promise<string> {
	const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
	const home = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/home.mustache')), { encoding: 'utf-8' });
	const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
	const messageEditForm = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/message_edit_form.mustache')), { encoding: 'utf-8' });
	const messageView = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/message_view.mustache')), { encoding: 'utf-8' });
	const newMessageForm = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/new_message_form.mustache')), { encoding: 'utf-8' });
	const homeSidebar = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/home_sidebar.mustache')), { encoding: 'utf-8' });
	const messagesList = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/messages_list.mustache')), { encoding: 'utf-8' });

	const content = mustache.render(home, {
		error:    err,
		domain:   Config.get("domain"),
		send:     i18n("send"),
		logout:   i18n("logout"),
		filesPlaceholder: i18n("filesPlaceholder"),
		newMessageText: i18n("newMessageText"),
		noMessagesText: i18n("noMessagesText"),
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
		scripts:  ["/public/pages/home/mobile.js"],
		manifest: "/public/manifest.json",
		styles:   styles,
	}, {
		content,
	});
}

async function renderBody(stack: Thread[], userID: number, message?: Message, editMessage?: boolean): Promise<string> {
	const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
	const home = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/home.mustache')), { encoding: 'utf-8' });
	const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
	const messageEditForm = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/message_edit_form.mustache')), { encoding: 'utf-8' });
	const messageView = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/message_view.mustache')), { encoding: 'utf-8' });
	const newMessageForm = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/new_message_form.mustache')), { encoding: 'utf-8' });
	const homeSidebar = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/home_sidebar.mustache')), { encoding: 'utf-8' });
	const messagesList = await readFile(resolve(join(Config.get('basedir'), 'templates/home/mobile/messages_list.mustache')), { encoding: 'utf-8' });

	const content = mustache.render(home, {
		stack:    stack,
		message:  message,
		userID:   userID,
		domain:   Config.get("domain"),
		send:     i18n("send"),
		logout:   i18n("logout"),
		filesPlaceholder: i18n("filesPlaceholder"),
		newMessageText: i18n("newMessageText"),
		noMessagesText: i18n("noMessagesText"),
		selectFiles: i18n("selectFiles"),
		search:   i18n("search"),
		delete:   i18n("delete"),
		edit:     i18n("edit"),
		editMessage: editMessage,
		limit: 14,
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
		scripts:  ["/public/pages/home/mobile.js"],
		manifest: "/public/manifest.json",
		styles:   styles,
	}, {
		content,
	});
}

export default renderMobile
