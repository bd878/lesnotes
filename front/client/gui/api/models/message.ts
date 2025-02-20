import file from './file';

const ns_in_ms = 10**6

export default function mapMessageFromProto(message) {
	if (!message) {
		return {
			ID: -1,
			createUTCNano: -1,
			updateUTCNano: -1,
			userID: 0,
			text: "",
			file: file(),
		}
	}

	let createUTCNano = 0
	if (message.create_utc_nano) {
		createUTCNano = new Date(Math.floor(message.create_utc_nano / ns_in_ms))
		createUTCNano = createUTCNano.toLocaleString()
	}

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
	if (message.file && message.file.id)
		res.file = file(message.file)

	return res
}