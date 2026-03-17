export interface Comment {
	messageID:     number;
	userID:        number;
	ID:            number;
	text:          string;
	createdAt:     string;
	updatedAt:     string;
}

const EmptyComment: Comment = Object.freeze({
	messageID:  0,
	userID:     0,
	ID:         0,
	text:       "",
	createdAt:  "",
	updatedAt:  "",
})

export default function mapCommentFromProto(comment?: any): Comment {
	if (!comment) {
		return EmptyComment
	}

	return {
		ID:          comment.id,
		messageID:   comment.message_id,
		userID:      comment.user_id,
		text:        comment.text,
		createdAt:   comment.createdAt,
		updatedAt:   comment.updatedAt,
	}
}

export { EmptyComment }
