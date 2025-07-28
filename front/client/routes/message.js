import Config from 'config';
import mustache from 'mustache';
import path from 'path';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

// - new message not static
// - public static message
async function renderMessage(ctx) {
	try {
		const filePath = resolve(path.join(Config.get('basedir'), 'templates/message.mustache'));
		const template = await readFile(filePath, { encoding: 'utf-8' });

		const token = ctx.cookies.get("token")
		console.log("token:", token)
		// public static
		if (is.empty(token)) {
			if (ctx.params.id) {
				const resp = await api.loadOneMessage(ctx.params.id)
				console.log('ctx.params.id, resp:', ctx.params.id, resp)

				ctx.body = mustache.render(template, {
					id: ctx.params.id,
					userId: ctx.params.userId,
					react: false,
					message: resp.message,
					styles: [
						"/public/styles.css",
					],
				});
			} else {
				ctx.body = mustache.render(template, {
					id: "",
					userId: ctx.params.userId || "",
					react: true,
					script: "/public/message.js",
					styles: [
						"/public/styles.css",
					],
				});				
			}
		} else if (!is.empty(token)) {
			ctx.body = mustache.render(template, {
				id: ctx.params.id || "",
				userId: ctx.params.userId || "",
				react: true,
				script: "/public/message.js",
				styles: [
					"/public/styles.css",
				],
			});
		}

		ctx.status = 200;
	} catch (err) {
		ctx.body = "<html>Pas de template</html>";
		ctx.status = 500;
		console.log("failed to return message template");
		throw Error(err);
	}
}

export default renderMessage;
