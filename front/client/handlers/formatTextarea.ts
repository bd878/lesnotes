import * as is from '../third_party/is';

async function formatTextarea(ctx, next) {
	console.log("--> formatTextarea")

	if (is.notEmpty(ctx.state.message) && is.notEmpty(ctx.state.message.text)) {
		ctx.state.message.text = ctx.state.message.text.replaceAll("\r\n", "&#13;&#10;")
	}
	if (is.notEmpty(ctx.state.thread) && is.notEmpty(ctx.state.thread.description)) {
		ctx.state.thread.description = ctx.state.thread.description.replaceAll("\r\n", "&#13;&#10;")
	}

	await next()

	console.log("<-- formatTextarea")
}

export default formatTextarea
