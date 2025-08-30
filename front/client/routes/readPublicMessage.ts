import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function readPublicMessage(ctx) {
	try {
		const id = parseInt(ctx.params.id, 10)
		const token = ctx.cookies.get("token")

		console.log(`[readPublicMessage]: token ${token} id ${id}`)

		const resp = await api.readMessageJson(token, 0, id)

		if (resp.error) {
			const filePath = resolve(join(Config.get('basedir'), 'templates/error.mustache'));
			const template = await readFile(filePath, { encoding: 'utf-8' });

			ctx.body = mustache.render(template, {
				code:     resp.code,
				explain:  resp.explain,
				styles:   ["/public/styles.css"],
			})

			ctx.status = resp.status
		} else {
			const filePath = resolve(join(Config.get('basedir'), 'templates/message.mustache'));
			const template = await readFile(filePath, { encoding: 'utf-8' });

			ctx.body = mustache.render(template, {
				id:       id,
				react:    false,
				message:  resp.message,
				styles:   ["/public/styles.css"],
			})

			ctx.status = 200;
		}
	} catch (err) {
		console.error("[readPublicMessage]: failed to return message template");
		throw Error(err);
	}
}

export default readPublicMessage;
