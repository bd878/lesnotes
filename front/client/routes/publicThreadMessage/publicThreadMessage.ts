import type { Message } from '../../api/models';
import Config from 'config';
import PublicThreadBuilder from '../publicThread/builder'
import mustache from 'mustache';
import * as is from '../../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

async function publicThreadMessage(ctx) {
	console.log("--> publicThreadMessage")

	const builder = new PublicThreadMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addMessagesList(ctx.params.threadName /* TODO: use from load_path, it is message now, thread required */, ctx.state.messages)
	await builder.addSearch()
	await builder.addMessageView(ctx.state.message)
	await builder.addSettings()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build(ctx.state.message)
	ctx.status = 200

	console.log("<-- publicThreadMessage")
}

class PublicThreadMessageBuilder extends PublicThreadBuilder {
	async addMessageView(message: Message) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/thread/mobile/message_view.mustache' : 'templates/thread/desktop/message_view.mustache'
		)), { encoding: 'utf-8' });

		this.messageView = mustache.render(template, {
			message: message,
		})
	}
}

export default publicThreadMessage
