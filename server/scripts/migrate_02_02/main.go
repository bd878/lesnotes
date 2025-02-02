package main

import (
  "flag"
  "os"
  "fmt"
  "context"
  "database/sql"

  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"

  _ "github.com/mattn/go-sqlite3"
  _ "github.com/bd878/gallery/server/messages/pkg/loadbalance"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/pkg/model"
  "github.com/bd878/gallery/server/api"
)

func init() {
  flag.Usage = func () {
    fmt.Printf("Usage: %s sourceFilePath\n", os.Args[0])
  }
}

const stmt = `
SELECT id, message, user_id FROM messages
`

func main() {
  flag.Parse()

  if flag.NArg() < 1 {
    flag.Usage()
    os.Exit(1)
  }

  log := logger.Default()

  sourceFilePath := flag.Arg(0)
  rpcAddr := "0.0.0.0:9001"

  controller := NewMessages(rpcAddr)

  log.Info("sourceFilePath=", sourceFilePath)

  pool, err := sql.Open("sqlite3", "file:" + sourceFilePath)
  if err != nil {
    panic(err)
  }

  rows, err := pool.QueryContext(context.Background(), stmt)
  if err != nil {
    panic(err)
  }
  defer rows.Close()

  for rows.Next() {
    var id int32
    var message string
    var userId int32
    if err := rows.Scan(
      &id,
      &message,
      &userId,
    ); err != nil {
      panic(err)
    }

    res, err := controller.SaveMessage(context.Background(), log, &model.SaveMessageParams{
      Message: &model.Message{
        Text: message,
        UserID: userId,
      },
    })
    if err != nil {
      log.Errorw("failed to save message", "id", id, "error", err)
      continue
    }
    log.Infow("message saved", "id", id, "create_utc_nano",
      res.CreateUTCNano, "update_utc_nano", res.UpdateUTCNano)
  }
}

type Messages struct {
  client  api.MessagesClient
  conn   *grpc.ClientConn
}

func NewMessages(rpcAddr string) *Messages {
  conn, err := grpc.Dial(
    fmt.Sprintf(
      "%s:///%s",
      "messages",
      rpcAddr,
    ),
    grpc.WithTransportCredentials(insecure.NewCredentials()),
  )
  if err != nil {
    panic(err)
  }

  client := api.NewMessagesClient(conn)

  return &Messages{client, conn}
}

func (s *Messages) Close() {
  if s.conn != nil {
    s.conn.Close()
  }
}

func (s *Messages) SaveMessage(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) (
  *model.SaveMessageResult, error,
) {
  res, err := s.client.SaveMessage(ctx, &api.SaveMessageRequest{
    Message: model.MessageToProto(params.Message),
  })
  if err != nil {
    log.Error("message", "client failed to save message")
    return nil, err 
  }

  return &model.SaveMessageResult{
    ID: res.Id,
    UpdateUTCNano: res.UpdateUtcNano,
    CreateUTCNano: res.CreateUtcNano,
  }, nil
}