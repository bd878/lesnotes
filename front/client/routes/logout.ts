import Config from 'config';
import api from '../api';
import i18n from '../i18n';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function logout(ctx) {
	const token = ctx.cookies.get("token")

	ctx.set({ "Cache-Control": "no-cache,max-age=0" })

	const resp = await api.authJson(token)
	if (resp.expired) {
		ctx.redirect('/login')
		ctx.status = 302
		return
	}

	const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
	const logout = await readFile(resolve(join(Config.get('basedir'), 'templates/logout.mustache')), { encoding: 'utf-8' });

	ctx.body = mustache.render(layout, {
		logout:   i18n("loading"),
		scripts:  ["/public/logout/logoutScript.js"],
	}, {
		content: logout,
	});

	ctx.status = 200;
}

export default logout;
