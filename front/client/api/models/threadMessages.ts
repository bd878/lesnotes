import type { Message } from './message'
import type { Paging } from './paging'
import paging from './paging'
import message from './message'

export interface ThreadMessages {
	messages: Message[];
	paging:   Paging;
	message:  Message;
	parentID: number;
	isRoot:   boolean;
// TODO: { parentID, isRoot } -> thread: Thread
}

const EmptyThreadMessages: ThreadMessages = Object.freeze({
	messages:   [],
	paging:     paging(),
	message:    message(),
	parentID:   0,
	isRoot:     true,
})

export default function mapThreadMessagesFromProto(messages: Message[], paging: Paging, message: Message, parentID: number, isRoot: boolean): ThreadMessages {
	return {
		messages:   messages,
		paging:     paging,
		message:    message,
// TODO: add thread: Thread, move parentID and isRoot under thread
		parentID:   parentID,
		isRoot:     isRoot,
	}
}

export { EmptyThreadMessages }
