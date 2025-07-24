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
			ctx.body = mustache.render(template, {
				token: resp.token,
				styles: [
					"/public/styles.css",
				],
			});
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
