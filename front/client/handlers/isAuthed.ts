function isAuthed(fun: (ctx: any, next: any) => Promise<any>): (ctx, next) => Promise<any> {
	return async function ifAuthed(ctx, next) {
		if (ctx.state.isAuthed) {
			await fun(ctx, next)
		} else {
			await next()
		}
	}
}

export default isAuthed