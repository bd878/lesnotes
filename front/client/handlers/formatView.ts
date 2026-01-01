import * as is from '../third_party/is';

async function formatView(ctx, next) {
	console.log("--> formatView")

	if (is.notEmpty(ctx.state.message) && is.notEmpty(ctx.state.message.text)) {
		ctx.state.message.text = ctx.state.message.text.replaceAll("\r\n", "<br />")
	}
	if (is.notEmpty(ctx.state.thread) && is.notEmpty(ctx.state.thread.description)) {
		ctx.state.thread.description = ctx.state.thread.description.replaceAll("\r\n", "<br />")
	}

	await next()

	console.log("<-- formatView")
}

export default formatView
