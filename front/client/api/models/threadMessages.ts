import type { Message } from './message'
import type { Paging } from './paging'
import type { Thread } from './thread'
import paging from './paging'
import message from './message'
import thread from './thread'

export interface ThreadMessages {
	messages: Message[];
	paging:   Paging;
	message?: Message;
	thread?:  Thread;
}

const EmptyThreadMessages: ThreadMessages = Object.freeze({
	messages:   [],
	paging:     paging(),
	message:    message(),
	thread:     thread(),
})

export default function mapThreadMessagesFromProto(messages: Message[], paging: Paging, message?: Message, thread?: Thread): ThreadMessages {
	return {
		messages:   messages,
		paging:     paging,
		message:    message,
		thread:     thread,
	}
}

export { EmptyThreadMessages }
