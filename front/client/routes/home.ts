import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import i18n from '../i18n';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function renderer(ctx) {
	const token = ctx.cookies.get("token")

	const resp = await api.authJson(token)
	if (resp.error.error || resp.expired) {
		ctx.redirect("/login")
		ctx.status = 302
		return
	}

	const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles.css')), { encoding: 'utf-8' });
	const home = await readFile(resolve(join(Config.get('basedir'), 'templates/home.mustache')), { encoding: 'utf-8' });
	const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
	const footer = await readFile(resolve(join(Config.get("basedir"), 'templates/footer.mustache')), { encoding: 'utf-8' });

	let browser = ""
	if (ctx.userAgent.isFirefox)
		browser = "firefox"
	else if (ctx.userAgent.isChrome)
		browser = "chrome"
	else if (ctx.userAgent.isSafari)
		browser = "safari"

	let mobile = false
	if (ctx.userAgent.isMobile)
		mobile = true

	ctx.body = mustache.render(layout, {
		scripts:  ["/public/home.js"],
		manifest: "/public/manifest.json",
		styles:   styles,
	}, {
		content: home,
		footer:  footer,
	});

	ctx.status = 200;
}

export default renderer;
