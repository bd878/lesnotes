export interface Thread {
	ID:           number;
	userID:       number;
	name:         string;
	private:      boolean;
	parentID:     number;
}

const EmptyThread = Object.freeze({
	ID:           0,
	userID:       0,
	name:         "",
	private:      true,
	parentID:     0,
})

export default function mapThreadFromProto(thread?: any): Thread {
	if (!thread)
		return EmptyThread

	const res = {
		ID:           thread.id,
		userID:       thread.user_id,
		name:         thread.name,
		parentID:     thread.parent_id,
		private:      thread.private
	}

	return res
}

export { EmptyThread }
