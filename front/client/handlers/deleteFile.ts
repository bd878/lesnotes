import * as is from '../third_party/is'
import api from '../api'

async function deleteFile(ctx) {
	console.log("--> deleteFile")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const response = await api.deleteFileJson(ctx.state.token, parseInt(form.id) || 0)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
	} else {
		ctx.redirect(ctx.router.url('files', {}, {query: ctx.query}))
	}

	console.log("<-- deleteFile")
}

export default deleteFile;
