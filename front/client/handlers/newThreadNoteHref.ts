import * as is from '../third_party/is'

async function newThreadNoteHref(ctx, next) {
	ctx.log.info("--> newThreadNoteHref")

	if (is.notEmpty(ctx.state.parentName)) {
		ctx.state.newNoteHref = "/" + ctx.state.parentName + "/new" + ctx.search
	} else if (is.notEmpty(ctx.state.message)) {
		ctx.state.newNoteHref = "/" + ctx.state.message.name + "/new" + ctx.search
	}

	await next()

	ctx.log.info("<-- newThreadNoteHref")
}

export default newThreadNoteHref
