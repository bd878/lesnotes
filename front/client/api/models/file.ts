interface File {
	ID:        number;
	name:      string;
}

const empty: File = {
	ID:   0,
	name: "",
}

export default function mapFileFromProto(file): File {
	if (!file)
		return empty

	return {
		ID:   file.id,
		name: file.name,
	}
}