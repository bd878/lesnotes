export interface File {
	ID:        number;
	name:      string;
}

export interface FileProto {
	id:        number;
	name:      string;
}

const EmptyFile: File = Object.freeze({
	ID:   0,
	name: "",
})

export default function mapFileFromProto(file?: FileProto): File {
	if (!file)
		return EmptyFile

	return {
		ID:   file.id,
		name: file.name,
	}
}

export { EmptyFile }
