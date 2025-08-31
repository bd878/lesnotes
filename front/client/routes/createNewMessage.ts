import Config from 'config';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function createNewMessage(ctx) {
	const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles.css')), { encoding: 'utf-8' });
	const template = await readFile(resolve(join(Config.get('basedir'), 'templates/new_message.mustache')), { encoding: 'utf-8' });

	ctx.body = mustache.render(template, {
		scripts:  ["/public/newScript.js"],
		styles:   styles,
	})

	ctx.status = 200;
}

export default createNewMessage;
