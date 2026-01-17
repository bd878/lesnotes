export interface File {
	ID:        number;
	name:      string;
	size:      number;
	mime:      string;
	private:   boolean;
}

export interface FileProto {
	id:        number;
	name:      string;
	size:      number;
	mime:      string;
	private:   boolean;
}

const EmptyFile: File = Object.freeze({
	ID:   0,
	name: "",
	size: 0,
	mime: "",
	private: true,
})

export default function mapFileFromProto(file?: FileProto): File {
	if (!file) {
		return EmptyFile
	}

	return {
		ID:      file.id,
		name:    file.name,
		size:    file.size,
		mime:    file.mime,
		private: file.private,
	}
}

export { EmptyFile }
