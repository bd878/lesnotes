import Builder from '../builder';

async function threadEdit(ctx) {
	console.log("--> threadEdit")

	const builder = new ThreadEditBuilder(ctx.userAgent.isMobile, ctx.state.lang)

	await builder.addSettings(ctx.state.lang, ctx.state.theme, ctx.state.fontSize)
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	console.log("<-- threadEdit")
}

class ThreadEditBuilder extends Builder {
	async build() {
		return "not implemented"
	}
}

export default threadEdit
