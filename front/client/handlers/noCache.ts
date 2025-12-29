async function noCache(ctx, next) {
	console.log("--> noCache")

	ctx.set({ 'Cache-Control': 'no-cache,max-age=0' })

	await next()

	console.log("<-- noCache")
}

export default noCache
