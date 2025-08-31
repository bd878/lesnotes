import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import i18n from '../i18n';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function main(ctx) {
	const token = ctx.cookies.get("token")

	const resp = await api.authJson(token)
	if (resp.error.error || resp.expired) {
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const main = await readFile(resolve(join(Config.get("basedir"), 'templates/main.mustache')), { encoding: 'utf-8' });
		const footer = await readFile(resolve(join(Config.get("basedir"), 'templates/footer.mustache')), { encoding: 'utf-8' });

		ctx.body = mustache.render(layout, {
			login:    i18n("login"),
			register: i18n("register"),
			scripts:  ["/public/mainScript.js"],
			styles:   ["/public/styles.css"],
		}, {
			content:  main,
			footer:   footer,
		});

		ctx.status = 200;

		return
	}

	ctx.redirect('/home')
	ctx.status = 301
}

export default main;
