import * as is from '../third_party/is'
import api from '../api'

async function updateMessage(ctx, next) {
	ctx.log.info("--> updateMessage")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	let fileIDs = []
	if (is.notEmpty(form.file_ids)) {
		if (is.array(form.file_ids)) {
			fileIDs = form.file_ids
		} else if (is.string(form.file_ids)) {
			fileIDs = JSON.parse(form.file_ids)
			if (!is.array(fileIDs)) {
				fileIDs = []
			}
		} else {
			fileIDs = [form.file_ids]
		}
	}

	fileIDs = fileIDs.map(id => parseInt(id) || 0).filter(is.notEmpty)

	const messageID = parseInt(form.id) || 0;

	// TODO: saga on backend
	response = await api.updateMessageJson(ctx.state.token, messageID, form.text, form.title, form.name, fileIDs)

	if (response.error.error) {
		ctx.log.error(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
		return
	}

	response = await api.updateThreadJson(ctx.state.token, messageID, "", "", form.name)
	if (response.error.error) {
		ctx.log.error(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
		return
	}

	const redirectUrl = form.redirectUrl

	if (is.notEmpty(redirectUrl)) {
		ctx.redirect(redirectUrl)
	} else {
		ctx.redirect(ctx.router.url("message", {idOrName: form.id}, {query: ctx.query}))
	}

	ctx.log.info("<-- updateMessage")
}

export default updateMessage