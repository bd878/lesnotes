import message from './message'
import user from './user'
import file, { EmptyFile } from './file'
import error from './error'

export default {
	message,
	user,
	file,
	error,
	EmptyFile,
}

export type { Message } from './message'
export type { File } from './file'
export type { Error } from './error'
export type { User } from './user'
