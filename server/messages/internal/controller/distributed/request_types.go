package distributed

type RequestType uint16

const (
	AppendRequest RequestType = iota
	UpdateRequest
	DeleteRequest
	PublishRequest
	PrivateRequest
	DeleteAllUserMessagesRequest
)
