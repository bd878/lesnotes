import Config from 'config';
import mustache from 'mustache';
import i18n from '../../i18n';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function createNewMessage(ctx) {
	const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
	const template = await readFile(resolve(join(Config.get('basedir'), 'templates/new_message.mustache')), { encoding: 'utf-8' });

	ctx.body = mustache.render(template, {
		scripts:  ["/public/newScript.js"],
		styles:   styles,
		send:     i18n('send'),
		file:     i18n("file"),
		text:     i18n("text"),
	})

	ctx.status = 200;
}

export default createNewMessage;
