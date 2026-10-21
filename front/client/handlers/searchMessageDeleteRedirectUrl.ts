async function searchMessageDeleteRedirectUrl(ctx, next) {
	ctx.log.info("--> deleteRedirectUrl")

	ctx.state.deleteRedirectUrl = "/search" + ctx.search

	await next()

	ctx.log.info("<-- deleteRedirectUrl")
}

export default searchMessageDeleteRedirectUrl
