package grpc

import (
  "io"
  "os"
  "time"
  "fmt"
  "context"
  "errors"
  "path/filepath"

  "github.com/bd878/gallery/server/utils"
  "github.com/bd878/gallery/server/api"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/files/pkg/model"
)

type Repository interface {
  SaveFile(ctx context.Context, log *logger.Logger, file *model.File) error
  ReadFile(ctx context.Context, log *logger.Logger, params *model.ReadFileParams) (*model.File, error)
}

type Handler struct {
  api.UnimplementedFilesServer
  repo       Repository
  dataPath   string
}

func New(repo Repository, dataPath string) *Handler {
  if err := os.MkdirAll(dataPath, 0755); err != nil {
    panic(err)
  }

  return &Handler{repo: repo, dataPath: dataPath}
}

func (h *Handler) ReadBatchFiles(ctx context.Context, req *api.ReadBatchFilesRequest) (
  *api.ReadBatchFilesResponse, error,
) {
  files := make(map[int32]*model.File, len(req.Ids))
  for _, id := range req.Ids {
    files[id] = &model.File{
      ID: id,
    }

    file, err := h.repo.ReadFile(ctx, logger.Default(), &model.ReadFileParams{ID: id})
    if err != nil {
      files[id].Error = "can not found file"
      logger.Errorw("failed to read file", "id", id, "error", err)
      continue
    }

    files[id] = file
  }

  return &api.ReadBatchFilesResponse{
    Files: model.MapFilesDictToProto(model.FileToProto, files),
  }, nil
}

func (h *Handler) ReadFileStream(params *api.ReadFileStreamRequest, stream api.Files_ReadFileStreamServer) error {
  file, err := h.repo.ReadFile(context.Background(), logger.Default(), &model.ReadFileParams{ID: params.Id})
  if err != nil {
    logger.Errorw("failed to read file", "id", params.Id, "error", err)
    return err
  }

  ff, err := os.Open(filepath.Join(h.dataPath, fmt.Sprintf("%d", file.ID)))
  if err != nil {
    logger.Errorw("failed to open file", "id", file.ID, "name", file.Name, "error", err)
    return err
  }

  err = stream.Send(&api.FileData{
    Data: &api.FileData_File{
      File: &api.File{
        Name:           file.Name,
        CreateUtcNano:  file.CreateUTCNano,
      },
    },
  })
  if err != nil {
    logger.Errorw("stream failed to send filedata", "id", file.ID, "error", err)
    return err
  }

  buffer := make([]byte, 1024*1024*20 /* 20 MB */)
  for {
    n, err := ff.Read(buffer)
    if err == io.EOF {
      break
    }
    if err != nil {
      logger.Errorw("filestream", "failed to read file data in buffer")
      return err
    }

    err = stream.Send(&api.FileData{
      Data: &api.FileData_Chunk{
        Chunk: buffer[:n],
      },
    })
    if err != nil {
      logger.Errorw("filestream", "failed to send chunk fil file server")
      return err
    }
  }

  return nil
}

func (h *Handler) SaveFileStream(stream api.Files_SaveFileStreamServer) error {
  meta, err := stream.Recv()
  if err != nil {
    logger.Errorw("save file stream failed to receive meta", "error", err)
    return err
  }

  file, ok := meta.Data.(*api.FileData_File)
  if !ok {
    logger.Errorw("wrong format", "send file data first, then chunk")
    return errors.New("wrong format: file meta expected")
  }

  id := utils.RandomID()
  timeCreated := time.Now().UnixNano()
  
  err = h.repo.SaveFile(context.Background(), logger.Default(), &model.File{
    ID:              id,
    Name:            file.File.Name,
    CreateUTCNano:   timeCreated,
  })
  if err != nil {
    logger.Errorw("failed to save file meta", "name", file.File.Name, "error", err)
    return err
  }

  ff, err := os.Create(filepath.Join(h.dataPath, fmt.Sprintf("%d", id)))
  if err != nil {
    logger.Errorw("failed to create file", "id", id, "error", err)
    return err
  }

  for {
    fileData, err := stream.Recv()
    if err == io.EOF {
      break
    }
    if err != nil {
      logger.Errorw("failed to receive next file chunk", "error", err)
      return err
    }

    chunk, ok := fileData.Data.(*api.FileData_Chunk)
    if !ok {
      logger.Errorw("wrong format", "file data chunk expected")
      return nil
    }

    n, err := ff.Write(chunk.Chunk)
    if err != nil {
      logger.Errorw("failed to write next file chunk in buffer", "error", err)
      return err
    }

    logger.Infow("write file", "id", id, "name", file.File.Name, "n", n)
  }

  return stream.SendAndClose(&api.SaveFileStreamResponse{
    File: &api.File{
      Id:               id,
      Name:             file.File.Name,
      CreateUtcNano:    timeCreated,
    },
  })
}