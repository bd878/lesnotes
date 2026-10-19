async function newNoteHref(ctx, next) {
	ctx.state.newNoteHref = "/home" + ctx.search
	await next()
}

export default newNoteHref;
