package distributed

type RequestType uint16

const (
  AppendRequest RequestType = 0
  UpdateRequest
  DeleteRequest
)
