import * as is from '../third_party/is';

async function formatMessage(ctx, next) {
	console.log("--> formatMessage")

	console.log(ctx.state.message)

	if (is.notEmpty(ctx.state.message) && is.notEmpty(ctx.state.message.text)) {
		ctx.state.message.text = ctx.state.message.text.replaceAll("\r\n", "&#13;&#10;")
	}

	await next()

	console.log("<-- formatMessage")
}

export default formatMessage
