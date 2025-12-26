import api from './index';
import models from './models';
import * as is from '../third_party/is'

async function readStackJson(token: string, threadID: number/*, lastMessageID: number*/, limit: number, offsets: Record<number, number> = {}) {
	let result = {
		error:       models.error(),
		stack:       [],
	}

	const path = await api.readPathJson(token, threadID)
	if (path.error.error) {
		result.error = path.error
		return result;
	}

	const ids = path.path.map(thread => thread.ID)
	ids.reverse()

	// len(threads) == len(path.path)
	const threads = [0 /* root thread */ , ...ids]

	path.path.push(JSON.parse(JSON.stringify(models.EmptyThread /* root thread */)))
	path.path.reverse()

	// threadID = 0 : threads = [0], path.path = [EmptyThread]
	for (
		let i = 0, isRoot = true, parentID = 0, threadID = 0 /* first is root : = 0 */, thread = JSON.parse(JSON.stringify(models.EmptyThread)) /* first is thread : EmptyThread */;
		i < threads.length;
		i++, isRoot = false, parentID = threads[i-1], threadID = threads[i], thread = path.path[i]
	) {
		const offset = offsets[threadID]

		let messages = { error: models.error(), messages: [], isLastPage: true, isFirstPage: true, count: 0, total: 0, offset: 0 }
		if (is.notUndef(offset)) {
			messages = await api.readMessagesJson(token, 0, threadID, 1 /* order */, limit, offset)
		} else {
			messages = await api.readMessagesJson(token, 0, threadID, 1 /* order */, limit, 0)
		}

		if (messages.error.error) {
			result.error = messages.error
			return result;
		}

		messages.messages.reverse()

		// TODO: move thread to api/models/thread.ts
		thread.isRoot      = isRoot /* parentID == 0 for root and primary child, how to distinguish? */
		thread.parentID    = parentID
		thread.isLastPage  = messages.isLastPage
		thread.isFirstPage = messages.isFirstPage
		thread.messages    = messages.messages
		thread.total       = messages.total
		thread.count       = messages.count
		thread.offset      = messages.offset

		result.stack.push(thread)
	}

	return result;
}

export default readStackJson;
