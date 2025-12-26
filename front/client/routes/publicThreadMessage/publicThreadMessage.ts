import type { Message } from '../../api/models';
import PublicThreadBuilder from '../publicThread/builder'
import * as is from '../../third_party/is';

async function publicThreadMessage(ctx) {
	console.log("--> publicThreadMessage")

	const builder = new PublicThreadMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addMessagesList(ctx.params.threadName /* TODO: use from load_path, it is message now, thread required */, ctx.state.stack)
	if (is.notEmpty(ctx.state.message)) {
		await builder.addMessageView(ctx.state.message)
	}
	await builder.addSettings()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	console.log("<-- publicThreadMessage")
}

class PublicThreadMessageBuilder extends PublicThreadBuilder {
	messageView = undefined
	async addMessageView(message: Message) {}
}

export default publicThreadMessage
