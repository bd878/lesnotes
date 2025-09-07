import Config from 'config';
import api from '../api';
import i18n from '../i18n';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function register(ctx) {
	const token = ctx.cookies.get("token")

	const resp = await api.authJson(token)
	if (resp.error.error || resp.expired) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const register = await readFile(resolve(join(Config.get("basedir"), 'templates/register.mustache')), { encoding: 'utf-8' });
		const footer = await readFile(resolve(join(Config.get("basedir"), 'templates/footer.mustache')), { encoding: 'utf-8' });

		ctx.body = mustache.render(layout, {
			username: i18n("username"),
			password: i18n("password"),
			register: i18n("register"),
			login:    i18n("login"),
			scripts:  ["/public/registerScript.js"],
			styles:   styles,
		}, {
			content: register,
			footer:  footer,
		});

		ctx.status = 200;

		return
	}

	ctx.redirect('/home')
	ctx.status = 302
}

export default register;
