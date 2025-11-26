// threadID is not in /search scope, we can't receive it here

export interface SearchMessage {
	ID:      number;
	userID:  number;
	text:    string;
	title:   string;
	name:    string;
	private: boolean;
}

const EmptySearchMessage: SearchMessage = Object.freeze({
	ID:      0,
	userID:  0,
	text:    "",
	title:   "",
	name:    "",
	private: true,
})

export default function mapMessageFromProto(message?: any): SearchMessage {
	if (!message)
		return EmptySearchMessage

	const res = {
		ID:        message.id,
		userID:    message.user_id,
		text:      message.text,
		name:      message.name,
		title:     message.title,
		private:   message.private,
	}

	return res
}

export { EmptySearchMessage }
