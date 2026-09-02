async function noCache(ctx, next) {
	ctx.log.info("--> noCache")

	ctx.set({ 'Cache-Control': 'no-cache,max-age=0' })

	await next()

	ctx.log.info("<-- noCache")
}

export default noCache
