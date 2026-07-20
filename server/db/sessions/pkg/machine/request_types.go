package machine

type RequestType uint16

const (
	AppendRequest RequestType = iota
	DeleteRequest
	DeleteUserSessionsRequest
)
