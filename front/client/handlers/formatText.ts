import * as is from '../third_party/is';

async function formatText(ctx, next) {
	console.log("--> formatText")

	if (is.notEmpty(ctx.state.message) && is.notEmpty(ctx.state.message.text)) {
		ctx.state.message.text = ctx.state.message.text.replaceAll("\r\n", "&#13;&#10;")
	}

	await next()

	console.log("<-- formatText")
}

export default formatText
