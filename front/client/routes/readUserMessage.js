import Config from 'config';
import mustache from 'mustache';
import path from 'path';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

async function readUserMessage(ctx) {
	try {
		const filePath = resolve(path.join(Config.get('basedir'), 'templates/message.mustache'));
		const template = await readFile(filePath, { encoding: 'utf-8' });

		const id = parseInt(ctx.params.id, 10)
		const user = parseInt(ctx.params.user, 10)
		const token = ctx.cookies.get("token")

		console.log(`[readUserMessage]: token ${token} user ${user} id ${id}`)

		const resp = await api.readMessageJson(token, user, id)

		if (resp.error) {
			ctx.body = "<html>" + resp.explain + "</html>"
			ctx.status = 500
			throw Error(resp.explain)
		} else {
			ctx.body = mustache.render(template, {
				id:       id,
				userId:   user,
				react:    false,
				message:  resp.message,
				styles:   ["/public/styles.css"],
			})

			ctx.status = 200;
		}
	} catch (err) {
		ctx.body = "<html>Pas de template</html>";
		ctx.status = 500;
		console.log("failed to return message template");
		throw Error(err);
	}
}

export default readUserMessage;
