import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import i18n from '../../i18n';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function main(ctx) {
	const token = ctx.cookies.get("token")

	const resp = await api.authJson(token)
	if (resp.error.error || resp.expired) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const main = await readFile(resolve(join(Config.get("basedir"), 'templates/main/desktop/main.mustache')), { encoding: 'utf-8' });
		const footer = await readFile(resolve(join(Config.get("basedir"), 'templates/footer.mustache')), { encoding: 'utf-8' });

		const _i18n = i18n(ctx.state.lang)

		ctx.body = mustache.render(layout, {
			login:    _i18n("login"),
			register: _i18n("register"),
			scripts:  ["/public/pages/main/mainScript.js"],
			styles:   styles,
		}, {
			content:  main,
			footer:  mustache.render(footer, {
				terms:    _i18n("terms"),
				contact:  _i18n("contact"),
				docs:     _i18n("docs"),
			}),
		});

		ctx.status = 200;

		return
	}

	ctx.redirect('/home')
	ctx.status = 302
}

export default main;
