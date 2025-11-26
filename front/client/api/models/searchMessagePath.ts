import type {SearchMessage} from './searchMessage'
import type {Message} from './message'

export interface SearchMessagePath {
	ID:      number;
	userID:  number;
	text:    string;
	title:   string;
	name:    string;
	private: boolean;
	path:    string[];
}
// thread paths, empty for root

const EmptySearchMessagePath: SearchMessagePath = Object.freeze({
	ID:      0,
	userID:  0,
	text:    "",
	title:   "",
	name:    "",
	private: true,
	path:    [],
})

export default function mapMessageFromProto(message?: SearchMessage, path: Message[] = []): SearchMessagePath {
	if (!message)
		return EmptySearchMessagePath

	const res = {
		ID:        message.ID,
		userID:    message.userID,
		text:      message.text,
		name:      message.name,
		title:     message.title,
		private:   message.private,
		path:      path.map(thread => thread.title), // TODO: title or text
	}

	return res
}

export { EmptySearchMessagePath }
