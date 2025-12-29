export interface ExternalThread {
	id:           number;
	user_id:      number;
	name:         string;
	parent_id:    string;
	private:      boolean;
	description:  string;
}

export interface Thread {
	ID:           number;
	userID:       number;
	name:         string;
	private:      boolean;
	parentID:     number;
	description:  string;
// TODO: add is_root, title, created_at, updated_at
}

const EmptyThread = Object.freeze({
	ID:           0,
	userID:       0,
	name:         "",
	private:      true,
	parentID:     0,
	description:  "",
})

export default function mapThreadFromProto(thread?: any): Thread {
	if (!thread)
		return EmptyThread

	const res = {
		ID:           thread.id,
		userID:       thread.user_id,
		name:         thread.name,
		parentID:     thread.parent_id,
		private:      thread.private,
		description:  thread.description,
	}

	return res
}

export { EmptyThread }
