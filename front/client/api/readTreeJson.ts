import type {Error, MessagesList} from './models'
import type {IDLimitOffset} from '../types'
import api from './api';
import models from './models';

interface ReadTreeResponse {
	error:    Error;
	list:     MessagesList;
}

async function readTreeJson(token: string, rootID: number, limit: number, offset: number, leaves: IDLimitOffset[]): Promise<ReadTreeResponse> {
	let result: ReadTreeResponse = {
		error:     models.error(),
		list:  models.EmptyMessagesList,
	}

	console.log("readTreeJson", "token", token, "root_id", rootID, "limit", limit, "offset", offset, "leaves", ...leaves)

	try {

		const [response, error] = await api('/messages/v2/read_tree', {
			method: "POST",
			body: {
				token: token,
				req:   {
					root:    rootID,
					limit:   limit,
					offset:  offset,
					leaves:  leaves,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			result.list = models.messagesList(response.list)
		}

	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default readTreeJson