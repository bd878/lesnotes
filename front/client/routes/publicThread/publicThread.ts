import Builder from '../builder';

async function publicThread(ctx) {
	console.log("--> publicThread")

	const builder = new ThreadBuilder(ctx.userAgent.isMobile, ctx.state.lang)

	await builder.addSettings(ctx.state.lang, ctx.state.theme, ctx.state.fontSize)
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	console.log("<-- publicThread")
}

class ThreadBuilder extends Builder {
	async build() {
		return "not implemented"
	}
}

export default publicThread
