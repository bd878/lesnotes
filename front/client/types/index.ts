import type { Message, File } from '../api/models'

interface FileWithMime extends File {
	isDocument: boolean;
	isImage:    boolean;
	isAudio:    boolean;
	isVideo:    boolean;
	isText:     boolean;
	isFile:     boolean;
}

export type { FileWithMime }
