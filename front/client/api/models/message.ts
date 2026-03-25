import type {File} from './file'
import type {Paging} from './paging'
import type {TranslationPreview} from './translationPreview'
import paging, {EmptyPaging} from './paging'
import file from './file';
import translationPreview from './translationPreview';
import * as is from '../../third_party/is'

const ns_in_ms = 10**6

export interface MessagesList extends Paging {
	messages: Message[]
}

export interface Message {
	ID:            number;
	createdAt:     string;
	updatedAt:     string;
	userID:        number;
	text:          string;
	name:          string;
	title:         string;
	count:         number;
	files:         File[];
	translations:  TranslationPreview[];
	messages:      MessagesList;
	private:       boolean;
}

const EmptyMessagesList: MessagesList = Object.freeze({
	...EmptyPaging,
	messages: [],
})

const EmptyMessage: Message = Object.freeze({
	ID: 0,
	createdAt: "",
	updatedAt: "",
	userID: 0,
	text: "",
	title: "",
	count: 0,
	name: "",
	files:  [],
	translations: [],
	messages: EmptyMessagesList,
	private: true,
})

function mapMessagesListFromProto(list?: any): MessagesList {
	if (!list) {
		return EmptyMessagesList
	}

	list.messages.reverse()

	const res = {
		...paging(list),
		messages: list.messages.map(mapMessageFromProto),
	}

	return res
}

export default function mapMessageFromProto(message?: any): Message {
	if (!message) {
		return EmptyMessage
	}

	const res = {
		ID:          message.id,
		createdAt:   message.createdAt,
		updatedAt:   message.updatedAt,
		userID:      message.user_id,
		text:        message.text,
		name:        message.name,
		title:       message.title,
		count:       message.count,
		private:     Boolean(message.private),
		messages:    EmptyMessagesList,
		files:       [],
		translations: [],
	}

	if (is.array(message.files)) {
		res.files = message.files.map(file)
	}

	if (is.array(message.translations)) {
		res.translations = message.translations.map(translationPreview)
	}

	if (is.notEmpty(message.messages)) {
		res.messages = mapMessagesListFromProto(message.messages)
	}

	return res
}

export { mapMessagesListFromProto as messagesList }
export { EmptyMessage }
export { EmptyMessagesList }
