export interface ExternalThread {
	id:           number;
	user_id:      number;
	name:         string;
	parent_id:    number;
	private:      boolean;
	description:  string;
	is_root:      boolean;
}

export interface Thread {
	ID:           number;
	userID:       number;
	name:         string;
	private:      boolean;
	parentID:     number;
	description:  string;
	isRoot:       boolean; // TODO: not implemented on server
// TODO: add is_root, title, created_at, updated_at
}

const EmptyThread = Object.freeze({
	ID:           0,
	userID:       0,
	name:         "",
	private:      true,
	parentID:     0,
	description:  "",
	isRoot:       false,
})

export default function mapThreadFromProto(thread?: ExternalThread): Thread {
	if (!thread) {
		return EmptyThread
	}

	const res = {
		ID:           thread.id,
		userID:       thread.user_id,
		name:         thread.name,
		parentID:     thread.parent_id,
		private:      thread.private,
		description:  thread.description,
		isRoot:       thread.is_root,
	}

	return res
}

export { EmptyThread }
