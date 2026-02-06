export interface Translation {
	messageID:      number;
	lang:           string;
	title:          string;
	text:           string;
	createdAt:      string;
	updatedAt:      string;
}

const EmptyTranslation: Translation = Object.freeze({
	messageID:  0,
	lang:       "",
	createdAt:  "",
	updatedAt:  "",
	title:      "",
	text:       "",
})

export default function mapTranslationFromProto(translation?: any): Translation {
	if (!translation) {
		return EmptyTranslation
	}

	let createdAt = "0"
	if (translation.created_at) {
		createdAt = new Date(Math.floor(translation.created_at / ns_in_ms)).toLocaleString()
	}

	let updatedAt = "0"
	if (translation.updated_at) {
		updatedAt = new Date(Math.floor(translation.updated_at / ns_in_ms)).toLocaleString()
	}

	const res = {
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		text:        translation.text,
		title:       translation.title,
		messageID:   translation.message,
		lang:        translation.lang,
	}

	return res
}

export { EmptyTranslation }
