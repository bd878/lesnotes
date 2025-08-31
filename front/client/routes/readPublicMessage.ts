import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function readPublicMessage(ctx) {
	const id = parseInt(ctx.params.id, 10)
	const token = ctx.cookies.get("token")

	console.log(`[readPublicMessage]: token ${token} id ${id}`)

	const resp = await api.readMessageJson(token, 0, id)

	if (resp.error.error) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles.css')), { encoding: 'utf-8' });
		const filePath = resolve(join(Config.get('basedir'), 'templates/error.mustache'));
		const template = await readFile(filePath, { encoding: 'utf-8' });

		ctx.body = mustache.render(template, {
			code:     resp.error.code,
			explain:  resp.error.explain,
			styles:   styles,
		})

		ctx.status = resp.error.status

		return
	}

	const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles.css')), { encoding: 'utf-8' });
	const filePath = resolve(join(Config.get('basedir'), 'templates/message.mustache'));
	const template = await readFile(filePath, { encoding: 'utf-8' });

	ctx.body = mustache.render(template, {
		id:       id,
		message:  resp.message,
		files:    resp.message.files,
		styles:   styles,
	})

	ctx.status = 200;
}

export default readPublicMessage;
