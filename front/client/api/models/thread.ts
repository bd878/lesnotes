import type {Message} from './message'

export interface Thread {
	ID:           number;
	userID:       number;
	messages:     Message[];
	centerID:     number;
	text:         string;
	title:        string;
	count:        number;
	name:         string;
	isLastPage:   boolean;
	isFirstPage:  boolean;
}

const EmptyThread = Object.freeze({
	ID:           0,
	userID:       0,
	name:         "",
	title:        "",
	text:         "",
	count:        0,
	messages:     [],
	centerID:     0,
	isLastPage:   false,
	isFirstPage:  false,
})

export default function mapThreadFromProto(thread?: any): Thread {
	if (!thread)
		return EmptyThread

	const res = {
		ID:           thread.id,
		userID:       thread.user_id,
		name:         thread.name,
		text:         thread.text,
		count:        thread.count,
		title:        thread.title,
		centerID:     thread.center_id || 0,
		messages:     thread.messages || [],
		isLastPage:   thread.is_last_page,
		isFirstPage:  thread.is_first_page,
	}

	return res
}

export { EmptyThread }
