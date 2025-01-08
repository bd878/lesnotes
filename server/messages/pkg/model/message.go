package model

type MessageId int

type FileId string

// This message handler passes to repository
type Message struct {
  Id MessageId          `json:"id"`
  CreateTime string     `json:"createtime"`
  UpdateUtcNano uint64  `json:"update_utc_nano"`
  UserId int            `json:"userid"`
  Value string          `json:"value"`
  FileName string       `json:"filename"`
  FileId FileId         `json:"fileid"`
  LogIndex uint64       `json:"logindex,omitempty"`
  LogTerm uint64        `json:"logterm,omitempty"`
}

type MessagesList struct {
  Messages   []*Message `json:"messages"`
  IsLastPage bool       `json:"islastpage"`
}

// Response to return to the client
type ServerResponse struct {
  Status string      `json:"status"`
  Description string `json:"description"`
}

type NewMessageServerResponse struct {
  ServerResponse
  Message Message `json:"message"`
}

type UpdateMessageServerResponse struct {
  ServerResponse
  Message Message `json:"message"`
}

type MessagesListServerResponse struct {
  ServerResponse
  Messages   []*Message `json:"messages"`
  IsLastPage bool       `json:"islastpage"`
}

const NullMsgId = MessageId(0)
