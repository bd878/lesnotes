import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function main(ctx) {
	const token = ctx.cookies.get("token")

	const resp = await api.authJson(token)
	if (resp.error.error || resp.expired) {
		const filePath = resolve(join(Config.get('basedir'), 'templates/main.mustache'));
		const template = await readFile(filePath, { encoding: 'utf-8' });

		ctx.body = mustache.render(template, {
			react:   false,
			scripts: ["/public/mainScript.js"],
			styles:  ["/public/styles.css"],
		});

		ctx.status = 200;

		return
	}

	ctx.redirect('/home')
	ctx.status = 301
}

export default main;
