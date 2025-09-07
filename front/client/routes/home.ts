import type { Message } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import i18n from '../i18n';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

/**
 * Renders home page
 * 
 * Example: /home?threads=[123,345]&id=111&limit=5&offset=10
 * threads - thread ids to load messages
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

	const limit = parseInt(ctx.query.limit) || 10
	const offset = parseInt(ctx.query.offset) || 0
	const threadID = parseInt(ctx.query.thread) || 0
	const id = parseInt(ctx.query.id) || 0

	const threads = await api.readMessagePathJson(token, threadID)
	if (threads.error.error) {
		ctx.body = await renderError("failed to load batch messages");
		ctx.status = 400;
		return;
	}

	const messages = await api.readMessagesJson(token, threadID, 0, limit, offset)
	if (messages.error.error) {
		ctx.body = await renderError("failed to load thread messages");
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

		ctx.body = await renderBody(reverse(threads.path), reverse(messages.messages), message.message)
		ctx.status = 200;

		return
	}

	ctx.body = await renderBody(reverse(threads.path), reverse(messages.messages))
	ctx.status = 200;

	return;
}

async function renderError(err: string): Promise<string> {
	const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles.css')), { encoding: 'utf-8' });
	const home = await readFile(resolve(join(Config.get('basedir'), 'templates/home.mustache')), { encoding: 'utf-8' });
	const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });

	return mustache.render(layout, {
		scripts:  ["/public/homeScript.js"],
		manifest: "/public/manifest.json",
		styles:   styles,
		error:    err,
		send:     i18n("send"),
		logout:   i18n("logout"),
		search:   i18n("search"),
		title_placeholder: i18n("title_placeholder"),
		text_placeholder:  i18n("text_placeholder"),
	}, {
		content: home,
	});
}

async function renderBody(threads: Message[], messages: Message[], message?: Message): Promise<string> {
	const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles.css')), { encoding: 'utf-8' });
	const home = await readFile(resolve(join(Config.get('basedir'), 'templates/home.mustache')), { encoding: 'utf-8' });
	const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });

	return mustache.render(layout, {
		scripts:  ["/public/homeScript.js"],
		manifest: "/public/manifest.json",
		styles:   styles,
		threads:  threads,
		messages: messages,
		message:  message,
		send:     i18n("send"),
		logout:   i18n("logout"),
		search:   i18n("search"),
		title_placeholder: i18n("title_placeholder"),
		text_placeholder:  i18n("text_placeholder"),
	}, {
		content: home,
	});
}

/* list methods */
function last(target: any[] = [], def: any = 0): any {
	if (target.length == 0)
		return def

	return target[target.length-1]
}

function head(target: any[] = [], def: any = 0): any {
	if (target.length == 0)
		return def

	return target[0]
}

function tail(target: any[] = [], def: any[] = []): any[] {
	if (target.length == 0)
		return def

	return target.slice(0, -1)
}

function reverse(target: any[] = [], def: any[] = []): any[] {
	if (target.length == 0)
		return def

	return target.reverse()
}

export default home;
