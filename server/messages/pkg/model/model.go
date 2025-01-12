package model

type SaveMessageParams struct {
  Message          *Message
}

type UpdateMessageParams struct {
  Message          *Message
}

type ReadUserMessagesParams struct {
  UserID            int32
  Limit             int32
  Offset            int32
  Ascending         bool
}

type AuthParams struct {
  Token             string
}

type PutParams struct {
  Message          *Message
}

type GetParams struct {
  UserID            int32
  Limit             int32
  Offset            int32
  Ascending         bool
}