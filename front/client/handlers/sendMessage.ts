import { setTimeout } from "node:timers/promises";
import * as is from '../third_party/is'
import api from '../api'
import models from '../api/models'

const limit = parseInt(LIMIT)

async function sendMessage(ctx) {
	// TODO: proxy send message to messages service, /send
	console.log("--> sendMessage")

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

	const response = await api.sendMessageJson(ctx.state.token, form.text, form.title, fileIDs, parseInt(form.thread) || 0, true)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
		return
	}

	await waitForThread(ctx, response.message.ID)

	const params = new URLSearchParams(ctx.query)
	params.set(form.thread, `${limit},0`)
	ctx.redirect(ctx.router.url('message', {idOrName: response.message.ID}, {query: params.toString()}))

	console.log("<-- sendMessage")
}

export default sendMessage;

async function waitForThread(ctx, threadID) {
	let response = { error: models.error() }
	let i = 0
	do {
		response = await api.readThreadJson(ctx.state.token, 0 /* me */, threadID)
		if (!response.error.error) {
			break
		}
		await setTimeout(500)
		console.log("waiting thread... ", i++, threadID)
	} while (response.error.error)
}