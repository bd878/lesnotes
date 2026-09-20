import { setTimeout } from "node:timers/promises";
import * as is from '../third_party/is'
import api from '../api'
import models from '../api/models'

async function deleteMessage(ctx) {
	ctx.log.info("--> deleteMessage")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const redirectUrl = form.deleteRedirectUrl

	const response = await api.deleteMessageJson(ctx.state.token, parseInt(form.id) || 0)

	if (response.error.error) {
		ctx.log.error(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
		return
	}

	await waitForThread(ctx, parseInt(form.id))

	if (is.notEmpty(redirectUrl)) {
		ctx.redirect(redirectUrl)
	} else {
		ctx.redirect(ctx.router.url('home', {idOrName: form.id}, {query: ctx.query}))
	}

	ctx.log.info("<-- deleteMessage")
}

export default deleteMessage;

async function waitForThread(ctx, threadID) {
	let response = { error: models.error() }
	let i = 0
	do {
		response = await api.readThreadJson(ctx.state.token, 0 /* me */, threadID)
		if (response.error.error && response.error.status == 404) {
			/* not found = deleted */
			ctx.log.info(JSON.stringify(response.error))
			break
		}

		await setTimeout(500)
		ctx.log.info(`waiting thread deleted... ${i}, threadID: ${threadID}`)
		i += 1
	} while (!response.error.error)
}
