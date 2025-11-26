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

	// len(threads) == len(centers) == len(path.path)
	const threads = [0 /* root thread */ , ...ids]
	const centers = [...ids, 0/* TODO: remove, lastMessageID is unused*/]

	path.path.push(JSON.parse(JSON.stringify(models.EmptyThread /* root thread */)))
	path.path.reverse()

	// threadID = 0 : threads = [0], centers = [lastMessageID], path.path = [EmptyThread]
	for (let i = 0; i < threads.length; i++) {
		const threadID = threads[i] /* first is root : = 0 */
		const centerID = centers[i] /* first is message : != 0 */
		const thread = path.path[i] /* first is thread : EmptyThread */

		const offset = offsets[threadID]

		let messages = { error: models.error(), messages: [], isLastPage: true, isFirstPage: true, count: 0, total: 0, offset: 0 }
		if (is.notUndef(offset))
			messages = await api.readMessagesJson(token, threadID, 1 /* order */, limit, offset)
		else if (is.notEmpty(centerID))
			messages = await api.readMessagesJson(token, threadID, 1 /* order */, limit, 0)
		else
			messages = await api.readMessagesJson(token, threadID, 1 /* order */, limit, 0)

		if (is.notEmpty(centerID))
			thread.centerID = centerID
		else
			thread.centerID = 0

		if (messages.error.error) {
			result.error = messages.error
			return result;
		}

		messages.messages.reverse()

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
