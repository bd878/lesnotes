const ns_in_ms = 10**6

export interface TranslationPreview {
	messageID:      number;
	lang:           string;
	title:          string;
	createdAt:      string;
	updatedAt:      string;
}

const EmptyTranslationPreview: TranslationPreview = Object.freeze({
	messageID:  0,
	lang:       "",
	createdAt:  "",
	updatedAt:  "",
	title:      "",
})

export default function mapTranslationPreviewFromProto(preview?: any): TranslationPreview {
	if (!preview) {
		return EmptyTranslationPreview
	}

	let createdAt = "0"
	if (preview.created_at) {
		createdAt = new Date(Math.floor(preview.created_at / ns_in_ms)).toLocaleString()
	}

	let updatedAt = "0"
	if (preview.updated_at) {
		updatedAt = new Date(Math.floor(preview.updated_at / ns_in_ms)).toLocaleString()
	}

	const res = {
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		title:       preview.title,
		messageID:   preview.message,
		lang:        preview.lang,
	}

	return res
}

export { EmptyTranslationPreview }
