async function parseThreadName(ctx, next) {
	console.log("--> parseThreadName")

	ctx.state.threadName = ctx.params.threadName

	await next()

	console.log("<-- parseThreadName")
}

export default parseThreadName
