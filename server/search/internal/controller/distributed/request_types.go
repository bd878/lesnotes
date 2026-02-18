package distributed

type RequestType uint16

const (
	AppendMessageRequest RequestType = iota
	UpdateMessageRequest
	DeleteMessageRequest
	PublishMessagesRequest
	PrivateMessagesRequest
	AppendThreadRequest
	UpdateThreadRequest
	DeleteThreadRequest
	ChangeThreadParentRequest
	PublishThreadRequest
	PrivateThreadRequest
	AppendFileRequest
	DeleteFileRequest
	PublishFileRequest
	PrivateFileRequest
	AppendTranslationRequest
	DeleteTranslationRequest
	UpdateTranslationRequest
)
