import file from './file';

const ns_in_ms = 10**6

const empty = {
	ID: 0,
	createUTCNano: 0,
	updateUTCNano: 0,
	userID: 0,
	text: "",
	fileID: 0,
	file: file(),
}

export default function mapMessageFromProto(message) {
	if (!message)
		return empty

	let createUTCNano = 0
	if (message.create_utc_nano) {
		createUTCNano = new Date(Math.floor(message.create_utc_nano / ns_in_ms))
		createUTCNano = createUTCNano.toLocaleString()
	}

// TODO: createUTCNano -> dateCreatedString
	let updateUTCNano = 0
	if (message.update_utc_nano) {
		updateUTCNano = new Date(Math.floor(message.update_utc_nano / ns_in_ms))
		updateUTCNano = updateUTCNano.toLocaleString()
	}

	const res = {
		ID: message.id,
		createUTCNano: createUTCNano,
		updateUTCNano: updateUTCNano,
		userID: message.user_id,
		text: message.text,
	}
	if (message.file && message.file_id) {
		res.fileID = message.file_id
		res.file = file(message.file)
	}

	return res
}