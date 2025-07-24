import Config from 'config';
import mustache from 'mustache';
import path from 'path';
import api from '../api';
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

async function authTelegram(ctx) {
	try {
		const resp = await api.validateTgAuthData(ctx.querystring)

		const filePath = resolve(path.join(Config.get('basedir'), 'templates/tg_auth.mustache'));
		const template = await readFile(filePath, { encoding: 'utf-8' });

		if (resp.ok) {
			var expireDate = new Date().getTime()
			const age = 1 * 24 * 60 * 60 * 1000 /* 1 day */
			expireDate += age

			ctx.cookies.set("token", resp.token, {
				maxAge: age,
				expires: new Date(expireDate),
				domain: Config.get("backend_url"),
				secure: true,
				overwrite: true,
			})

			ctx.redirect('/')
			ctx.status = 301
		} else {
			ctx.body = mustache.render(template, {
				error: resp.error,
				explain: resp.explain,
				styles: [
					"/public/styles.css",
				],
			});
		}

		ctx.status = 200;
	} catch (err) {
		ctx.body = "<html>Pas de template</html>";
		ctx.status = 500;
		console.log("failed to auth telegram");
		throw Error(err);
	}
}

export default authTelegram;
