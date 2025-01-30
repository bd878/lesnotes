package model

type Message struct {
  ID                  int32         `json:"id"`
  CreateUTCNano       int64         `json:"create_utc_nano"`
  UpdateUTCNano       int64         `json:"update_utc_nano"`
  UserID              int32         `json:"user_id"`
  FileID              int32         `json:"file_id"`
  Text                string        `json:"text"`
}

type MessagesList struct {
  Messages            []*Message    `json:"messages"`
  IsLastPage          bool          `json:"is_last_page"`
}

type ServerResponse struct {
  Status              string        `json:"status"`
  Description         string        `json:"description"`
}

type NewMessageServerResponse struct {
  ServerResponse
  Message             Message       `json:"message"`
}

type UpdateMessageServerResponse struct {
  ServerResponse
  ID                  int32         `json:"id"`
  UpdateUTCNano       int64         `json:"update_utc_nano"`
}

type DeleteMessageServerResponse struct {
  ServerResponse
  ID                  int32         `json:"id"`
}

type MessagesListServerResponse struct {
  ServerResponse
  Messages            []*Message    `json:"messages"`
  IsLastPage          bool          `json:"is_last_page"`
}
