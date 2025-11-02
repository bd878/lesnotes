async function getFontSize(ctx, next) {
	// may be mutated by getMe
	console.log("--> getFontSize")
	await next()
	console.log("<-- getFontSize")
}

export default getFontSize
