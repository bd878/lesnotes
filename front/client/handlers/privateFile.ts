import * as is from '../third_party/is'
import api from '../api'

async function privateFile(ctx) {
	console.log("--> privateFile")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const response = await api.privateFileJson(ctx.state.token, parseInt(form.id) || 0)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
	} else {
		ctx.redirect(ctx.router.url('files', {}, {query: ctx.query}))
	}

	console.log("<-- privateFile")
}

export default privateFile;
