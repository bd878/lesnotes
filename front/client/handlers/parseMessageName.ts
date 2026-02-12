async function parseMessageName(ctx, next) {
	console.log("--> parseMessageName")

	ctx.state.messageName = ctx.params.messageName

	await next()

	console.log("<-- parseMessageName")
}

export default parseMessageName
