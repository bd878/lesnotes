async function parseMessageID(ctx, next) {
	console.log("--> parseMessageID")

	ctx.state.messageID = parseInt(ctx.params.id) || 0

	await next()

	console.log("<-- parseMessageID")
}

export default parseMessageID
