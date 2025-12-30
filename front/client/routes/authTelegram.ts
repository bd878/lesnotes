import Config from 'config';
import mustache from 'mustache';
import path from 'path';
import api from '../api';
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

async function authTelegram(ctx) {
	console.log("--> authTelegram")

	const resp = await api.validateTgAuthData(ctx.querystring)

	const filePath = resolve(path.join(Config.get('basedir'), 'templates/tg_auth.mustache'));
	const template = await readFile(filePath, { encoding: 'utf-8' });

	if (!resp.error.error) {
		var expireDate = new Date().getTime()
		const age = 1 * 24 * 60 * 60 * 1000 /* 1 day */
		expireDate += age

		ctx.cookies.set("token", resp.token, {
			maxAge: age,
			expires: new Date(expireDate),
			domain: Config.get("backend_url"),
			secure: false,
			overwrite: true,
		})

		ctx.redirect('/home')
		ctx.status = 301
	} else {
		ctx.body = mustache.render(template, {
			error: resp.error.error,
			explain: resp.error.explain,
			styles: [
				"/public/styles/styles.css",
			],
		});
		ctx.status = 200;
	}

	console.log("<-- authTelegram")
}

export default authTelegram;
