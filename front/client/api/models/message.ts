import type {File} from './file'
import file from './file';
import * as is from '../../third_party/is'

const ns_in_ms = 10**6

export interface Message {
	ID:            number;
	createUTCNano: string;
	updateUTCNano: string;
	userID:        number;
	text:          string;
	name:          string;
	files:         File[];
	threadID:      number;
	private:       boolean;
}

const empty: Message = {
	ID: 0,
	createUTCNano: "",
	updateUTCNano: "",
	userID: 0,
	text: "",
	name: "",
	files:  [],
	threadID: 0,
	private: true,
}

export default function mapMessageFromProto(message?: any): Message {
	if (!message)
		return empty

	let createUTCNano = "0"
	if (message.create_utc_nano) {
		createUTCNano = new Date(Math.floor(message.create_utc_nano / ns_in_ms)).toLocaleString()
	}

	let updateUTCNano = "0"
	if (message.update_utc_nano) {
		updateUTCNano = new Date(Math.floor(message.update_utc_nano / ns_in_ms)).toLocaleString()
	}

	const res = {
		ID: message.id,
		createUTCNano: createUTCNano,
		updateUTCNano: updateUTCNano,
		userID: message.user_id,
		text: message.text,
		name: "",
		private: Boolean(message.private),
		threadID: message.thread_id,
		files: [],
	}
	if (is.array(message.files)) {
		res.files = message.files.map(file)
	}

	return res
}