import * as is from '../third_party/is'

async function redirectSearch(ctx, next) {
	ctx.log.info("--> redirectSearch")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const params = new URLSearchParams(ctx.query)
	params.set("q", form.query)
	console.log("params:", params.toString())
	ctx.redirect(ctx.router.url('search', {}, {query: params.toString()}))
	ctx.status = 302

	ctx.log.info("<-- redirectSearch")
}

export default redirectSearch
