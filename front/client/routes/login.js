import Config from 'config';
import path from 'path';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

async function renderer(ctx) {
	try {
		const mainPath = resolve(path.join(Config.get('basedir'), 'templates/main.mustache'));
		const loginPath = resolve(path.join(Config.get('basedir'), 'templates/login.mustache'));
		const mainTemplate = await readFile(mainPath, { encoding: 'utf-8' });
		const loginTemplate = await readFile(loginPath, { encoding: 'utf-8' });

		ctx.body = mustache.render(
			mainTemplate,
			{
				script: "/public/login.js",
				styles: ["/public/styles.css"],
				form: "test",
				button_text: "Отправить",
			},
			{
				content: loginTemplate,
			},
		);
		ctx.status = 200;
	} catch (err) {
		ctx.body = "<html>Pas de template</html>";
		ctx.status = 500;
		console.log("failed to return index template");
		throw Error(err);
	}
}

export default renderer;
