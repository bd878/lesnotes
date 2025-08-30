import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function readUserMessage(ctx) {
	const id = parseInt(ctx.params.id, 10)
	const user = parseInt(ctx.params.user, 10)
	const token = ctx.cookies.get("token")

	console.log(`[readUserMessage]: token ${token} user ${user} id ${id}`)

	const resp = await api.readMessageJson(token, user, id)

	if (resp.error.error) {
		const filePath = resolve(join(Config.get('basedir'), 'templates/error.mustache'));
		const template = await readFile(filePath, { encoding: 'utf-8' });

		ctx.body = mustache.render(template, {
			code:     resp.error.code,
			explain:  resp.error.explain,
			styles:   ["/public/styles.css"],
		})

		ctx.status = resp.error.status

		return
	}

	const filePath = resolve(join(Config.get('basedir'), 'templates/message.mustache'));
	const template = await readFile(filePath, { encoding: 'utf-8' });

	ctx.body = mustache.render(template, {
		id:       id,
		user:     user,
		message:  resp.message,
		files:    resp.message.files,
		styles:   ["/public/styles.css"],
	})

	ctx.status = 200;
}

export default readUserMessage;
