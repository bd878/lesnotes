import message, { EmptyMessage } from './message'
import translation, { EmptyTranslation } from './translation';
import searchMessage, { EmptySearchMessage } from './searchMessage'
import searchMessagePath, { EmptySearchMessagePath } from './searchMessagePath'
import threadMessages, { EmptyThreadMessages } from './threadMessages'
import thread, { EmptyThread } from './thread'
import paging, { EmptyPaging, unwrapPaging } from './paging'
import user from './user'
import file, { EmptyFile } from './file'
import error from './error'

export default {
	message,
	searchMessage,
	searchMessagePath,
	threadMessages,
	translation,
	thread,
	paging,
	user,
	file,
	error,
	EmptyFile,
	EmptyPaging,
	EmptyThread,
	EmptyMessage,
	EmptyTranslation,
	EmptyThreadMessages,
	EmptySearchMessage,
	EmptySearchMessagePath,

	unwrapPaging,
}

export type { Thread } from './thread'
export type { Message } from './message'
export type { Translation } from './translation'
export type { File } from './file'
export type { Error } from './error'
export type { User } from './user'
export type { Paging } from './paging'
export type { ThreadMessages } from './threadMessages'
export type { SearchMessage } from './searchMessage'
export type { SearchMessagePath } from './searchMessagePath'
