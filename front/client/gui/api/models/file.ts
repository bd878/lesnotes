const empty = {
	ID: 0,
	name: "",
}
export default function mapFileFromProto(file) {
	if (!file)
		return empty

	return {
		ID: file.id,
		name: file.name,
	}
}