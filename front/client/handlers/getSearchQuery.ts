import * as is from '../third_party/is';

async function getSearchQuery(ctx, next) {
	console.log('--> getSearchQuery')

	ctx.state.query = ctx.query.query
	await next()

	console.log('<-- getSearchQuery')
}

export default getSearchQuery
