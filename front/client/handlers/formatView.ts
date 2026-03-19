import * as is from '../third_party/is';

async function formatView(ctx, next) {
	console.log("--> formatView")

	if (is.notEmpty(ctx.state.message) && is.notEmpty(ctx.state.message.text)) {
		ctx.state.message.text = ctx.state.message.text.replaceAll("\r\n", "<br />")
	}
	if (is.notEmpty(ctx.state.thread) && is.notEmpty(ctx.state.thread.description)) {
		ctx.state.thread.description = ctx.state.thread.description.replaceAll("\r\n", "<br />")
	}
	if (is.notEmpty(ctx.state.translation) && is.notEmpty(ctx.state.translation.text)) {
		ctx.state.translation.text = ctx.state.translation.text.replaceAll("\r\n", "<br />")
	}
	if (is.notEmpty(ctx.state.comments) && is.array(ctx.state.comments)) {
		ctx.state.comments = ctx.state.comments.map(comment => { comment.text = comment.text.replaceAll("\r\n", "<br />").replaceAll("\r", "<br />"); return comment })
	}

	await next()

	console.log("<-- formatView")
}

export default formatView
