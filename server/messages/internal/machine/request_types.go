package machine

type RequestType uint16

const (
	AppendRequest RequestType = iota
	UpdateRequest
	DeleteRequest
	PublishRequest
	PrivateRequest
	DeleteUserMessagesRequest
	AppendTranslationRequest
	UpdateTranslationRequest
	DeleteTranslationRequest

	AppendCommentRequest
	UpdateCommentRequest
	DeleteCommentRequest
	DeleteMessageCommentsRequest
)
