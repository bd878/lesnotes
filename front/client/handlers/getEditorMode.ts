import * as is from '../third_party/is';

async function getEditorMode(ctx, next) {
	console.log("--> getEditorMode")

	if (is.notEmpty(ctx.state.message)) {
		if (ctx.query.edit) {
			ctx.state.editorMode = "edit"
		} else {
			ctx.state.editorMode = "view"
		}
	} else {
		ctx.state.editorMode = "new-message"
	}

	await next()

	console.log("<-- getEditorMode")
}

export default getEditorMode
