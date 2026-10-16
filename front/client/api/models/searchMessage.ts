import type {Paging} from './paging'
import paging, {EmptyPaging} from './paging'
// threadID is not in /search scope, we can't receive it here

export interface SearchMessage {
	ID:      number;
	userID:  number;
	text:    string;
	title:   string;
	name:    string;
	private: boolean;
}

export interface SearchMessagesList extends Paging {
	messages: SearchMessage[]
}

const EmptySearchMessage: SearchMessage = Object.freeze({
	ID:      0,
	userID:  0,
	text:    "",
	title:   "",
	name:    "",
	private: true,
})

const EmptySearchMessagesList: SearchMessagesList = Object.freeze({
	...EmptyPaging,
	messages: [],
})

function mapSearchMessagesListFromProto(list?: any): SearchMessagesList {
	if (!list) {
		return EmptySearchMessagesList
	}

	return {
		...paging(list),
		messages: list.list.map(mapMessageFromProto),
	}
}

export default function mapMessageFromProto(message?: any): SearchMessage {
	if (!message) {
		return EmptySearchMessage
	}

	const res = {
		ID:        message.id,
		userID:    message.user_id,
		text:      message.text,
		name:      message.name,
		title:     message.title,
		private:   message.private,
	}

	return res
}

export { EmptySearchMessage }
export { mapSearchMessagesListFromProto as searchMessagesList }
export { EmptySearchMessagesList }
