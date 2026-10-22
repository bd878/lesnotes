async function newNoteHref(ctx, next) {
	if (ctx.userAgent.isMobile) {
		ctx.state.newNoteHref = "/new" + ctx.search
	} else {
		ctx.state.newNoteHref = "/home" + ctx.search
	}
	await next()
}

export default newNoteHref;
