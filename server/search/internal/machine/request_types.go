package machine

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
	DeleteFilesRequest
	PublishFilesRequest
	PrivateFilesRequest
	AppendTranslationRequest
	DeleteTranslationRequest
	UpdateTranslationRequest
)
