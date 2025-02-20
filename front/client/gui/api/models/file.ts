export default function mapFileFromProto(file) {
	if (!file) {
		return {
			ID: -1,
			name: "",
		}
	}

	return {
		ID: file.id,
		name: file.name,
	}
}