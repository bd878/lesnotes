import Config from 'config';
import mustache from 'mustache';
import path from 'path';
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

async function renderer(ctx) {
	try {
		const filePath = resolve(path.join(Config.get('basedir'), 'templates/index.mustache'));
		const template = await readFile(filePath, { encoding: 'utf-8' });

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

		ctx.body = mustache.render(template, {
			script: "/public/home.js",
			manifest: "/public/manifest.json",
			browser: browser,
			mobile: mobile,
			styles: [
				"/public/styles.css",
			],
		});
		ctx.status = 200;
	} catch (err) {
		ctx.body = "<html>Pas de template</html>";
		ctx.status = 500;
		console.log("failed to return index template");
		throw Error(err);
	}
}

export default renderer;
