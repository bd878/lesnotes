import * as is from '../third_party/is';

// TODO: remove, same logic is in ctx.search
async function getQuery(ctx, next) {
	console.log("--> getQuery")
	ctx.state.query = is.notEmpty(ctx.querystring) ? "?" + ctx.querystring : ""
	await next()
	console.log("<-- getQuery")
}

export default getQuery
