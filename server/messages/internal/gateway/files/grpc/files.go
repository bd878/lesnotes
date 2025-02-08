package grpc

import (
  "context"
  "io"

  "google.golang.org/grpc"

  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/pkg/model"
  filesmodel "github.com/bd878/gallery/server/files/pkg/model"
  "github.com/bd878/gallery/server/internal/grpcutil"
)

type Gateway struct {
  filesAddr string
  client    api.FilesClient
  conn      *grpc.ClientConn
}

func New(filesAddr string) *Gateway {
  conn, err := grpcutil.ServiceConnection(context.Background(), filesAddr)
  if err != nil {
    panic(err)
  }

  client := api.NewFilesClient(conn)

  return &Gateway{filesAddr, client, conn}
}

func (g *Gateway) SaveFileStream(ctx context.Context, log *logger.Logger, fileStream io.Reader, params *model.SaveFileParams) (
  *model.SaveFileResult, error,
) {
  stream, err := g.client.SaveFileStream(ctx)
  if err != nil {
    log.Error("message", "client failed to obtain file stream")
    return nil, err
  }

  err = stream.Send(&api.FileData{
    Data: &api.FileData_File{
      File: &api.File{
        Name: params.Name,
      },
    },
  })
  if err != nil {
    log.Error("message", "failed to save file meta")
    return nil, err
  }

  buffer := make([]byte, 1024)
  for {
    n, err := fileStream.Read(buffer)
    if err == io.EOF {
      break
    }
    if err != nil {
      log.Error("filestream", "failed to read file data in buffer")
      return nil, err
    }

    err = stream.Send(&api.FileData{
      Data: &api.FileData_Chunk{
        Chunk: buffer[:n],
      },
    })
    if err != nil {
      log.Error("filestream", "failed to send chunk fil file server")
      return nil, err
    }
  }

  res, err := stream.CloseAndRecv()
  if err != nil {
    log.Error("filestream", "failed to close and recv result")
    return nil, err
  }

  return &model.SaveFileResult{
    ID:              res.File.Id,
    CreateUTCNano:   res.File.CreateUtcNano,
  }, nil
}

func (g *Gateway) ReadBatchFiles(ctx context.Context, log *logger.Logger, params *model.ReadBatchFilesParams) (
  *model.ReadBatchFilesResult, error,
) {
  batch, err := g.client.ReadBatchFiles(ctx, &api.ReadBatchFilesRequest{
    Ids: params.IDs,
  })
  if err != nil {
    log.Error("message", "client failed to read batch files")
    return nil, err
  }

  return &model.ReadBatchFilesResult{
    Files: filesmodel.MapFilesDictFromProto(filesmodel.FileFromProto, batch.Files),
  }, nil
}