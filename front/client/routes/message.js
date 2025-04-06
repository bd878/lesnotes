import Config from 'config';
import mustache from 'mustache';
import path from 'path';
import api from '../api';
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

async function renderMessage(ctx) {
	try {
		const filePath = resolve(path.join(Config.get('basedir'), 'templates/message.mustache'));
		const template = await readFile(filePath, { encoding: 'utf-8' });

		if (ctx.params.id) {
			try {
				const resp = await api.loadOneMessage(ctx.params.id)
				console.log('ctx.params.id, resp:', ctx.params.id, resp)
			} catch (e) {
				console.log("error:", resp)
			}
		}

		ctx.body = mustache.render(template, {
			id: ctx.params.id || "",
			react: true,
			script: "/public/message.js",
			styles: [
				"/public/styles.css",
			],
		});
		ctx.status = 200;
	} catch (err) {
		ctx.body = "<html>Pas de template</html>";
		ctx.status = 500;
		console.log("failed to return message template");
		throw Error(err);
	}
}

export default renderMessage;
