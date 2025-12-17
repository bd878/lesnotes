import * as is from '../third_party/is';

async function getSearchForm(ctx, next) {
	console.log('--> getSearchForm')

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	ctx.state.query = form.search

	const search = new URLSearchParams(ctx.search)
	search.set("query", form.search)
	ctx.search = "?" + search.toString()

	await next()

	console.log('<-- getSearchForm')
}

export default getSearchForm
