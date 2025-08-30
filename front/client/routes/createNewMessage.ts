import Config from 'config';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function createNewMessage(ctx) {
	const filePath = resolve(join(Config.get('basedir'), 'templates/new_message.mustache'));
	const template = await readFile(filePath, { encoding: 'utf-8' });

	ctx.body = mustache.render(template, {
		styles:   ["/public/styles.css"],
		scripts:  ["/public/newScript.js"],
	})

	ctx.status = 200;
}

export default createNewMessage;
