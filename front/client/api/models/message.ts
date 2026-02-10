import type {File} from './file'
import type {Translation} from './translation'
import file from './file';
import translation from './translation';
import * as is from '../../third_party/is'

const ns_in_ms = 10**6

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
	translations:  Translation[];
	private:       boolean;
}

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
	private: true,
})

export default function mapMessageFromProto(message?: any): Message {
	if (!message) {
		return EmptyMessage
	}

	let createdAt = "0"
	if (message.create_utc_nano) {
		createdAt = new Date(Math.floor(message.create_utc_nano / ns_in_ms)).toLocaleString()
	}

	let updatedAt = "0"
	if (message.update_utc_nano) {
		updatedAt = new Date(Math.floor(message.update_utc_nano / ns_in_ms)).toLocaleString()
	}

	const res = {
		ID:          message.id,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		userID:      message.user_id,
		text:        message.text,
		name:        message.name,
		title:       message.title,
		count:       message.count,
		private:     Boolean(message.private),
		files:       [],
		translations: [],
	}

	if (is.array(message.files)) {
		res.files = message.files.map(file)
	}

	if (is.array(message.translations)) {
		res.translations = message.translations.map(translation)
	}

	return res
}

export { EmptyMessage }
