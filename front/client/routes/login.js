import Config from 'config';
import path from 'path';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

async function renderer(ctx) {
	try {
		const templatePath = resolve(path.join(Config.get('basedir'), 'templates/login.mustache'));
		const template = await readFile(templatePath, { encoding: 'utf-8' });

		ctx.body = mustache.render(
			template,
			{
				backend_url: Config.get("backend_url"),
				script: "/public/login.js",
				styles: ["/public/styles.css"],
			},
		);
		ctx.status = 200;
	} catch (err) {
		ctx.body = "<html>Pas de template</html>";
		ctx.status = 500;
		console.log("failed to return login template");
		throw Error(err);
	}
}

export default renderer;
