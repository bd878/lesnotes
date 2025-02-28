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
      ID:     id,
      UserID: req.UserId,
    }

    file, err := h.repo.ReadFile(ctx, logger.Default(), &model.ReadFileParams{ID: id, UserID: req.UserId})
    if err != nil {
      files[id].Error = "can not found file"
      logger.Errorw("failed to read file", "user_id", req.UserId, "id", id, "error", err)
      continue
    }

    files[id] = file
  }

  return &api.ReadBatchFilesResponse{
    Files: model.MapFilesDictToProto(model.FileToProto, files),
  }, nil
}

func (h *Handler) ReadFileStream(params *api.ReadFileStreamRequest, stream api.Files_ReadFileStreamServer) error {
  file, err := h.repo.ReadFile(context.Background(), logger.Default(), &model.ReadFileParams{ID: params.Id, UserID: params.UserId})
  if err != nil {
    logger.Errorw("failed to read file", "user_id", params.UserId, "id", params.Id, "error", err)
    return err
  }

  ff, err := os.Open(filepath.Join(h.dataPath, fmt.Sprintf("%d/%d", file.UserID, file.ID)))
  if err != nil {
    logger.Errorw("failed to open file", "user_id", file.UserID, "id", file.ID, "name", file.Name, "error", err)
    return err
  }

  var size int64
  stat, err := ff.Stat()
  if err != nil {
    logger.Errorw("cannot stat file", "error", err, "name", file.Name, "id", file.ID)
  } else {
    size = stat.Size()
  }

  err = stream.Send(&api.FileData{
    Data: &api.FileData_File{
      File: &api.File{
        Id:             file.ID,
        UserId:         file.UserID,
        Name:           file.Name,
        CreateUtcNano:  file.CreateUTCNano,
        Size:           size,
      },
    },
  })
  if err != nil {
    logger.Errorw("stream failed to send filedata", "user_id", file.UserID, "id", file.ID, "error", err)
    return err
  }

  buffer := make([]byte, 1024*1024*50 /* 50 MB */)
  for {
    n, err := ff.Read(buffer)
    if err == io.EOF {
      break
    }
    if err != nil {
      logger.Errorw("failed to read file data in buffer", "error", err)
      return err
    }

    err = stream.Send(&api.FileData{
      Data: &api.FileData_Chunk{
        Chunk: buffer[:n],
      },
    })
    if err != nil {
      logger.Errorw("failed to send chunk fil file server", "error", err)
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
    logger.Errorw("send file data first, then chunk", "error", "wrong format")
    return errors.New("wrong format: file meta expected")
  }

  id := utils.RandomID()
  timeCreated := time.Now().UnixNano()
  
  err = h.repo.SaveFile(context.Background(), logger.Default(), &model.File{
    ID:              id,
    UserID:          file.File.UserId,
    Name:            file.File.Name,
    CreateUTCNano:   timeCreated,
  })
  if err != nil {
    logger.Errorw("failed to save file meta", "user_id", file.File.UserId, "name", file.File.Name, "error", err)
    return err
  }

  userDir := filepath.Join(h.dataPath, fmt.Sprintf("%d", file.File.UserId))
  err = os.MkdirAll(userDir, 0755)
  if err != nil {
    logger.Errorw("cannot create user files dir", "user_id", file.File.UserId, "error", err)
    return err
  }

  ff, err := os.Create(filepath.Join(h.dataPath, fmt.Sprintf("%d/%d", file.File.UserId, id)))
  if err != nil {
    logger.Errorw("failed to create file", "user_id", file.File.UserId, "id", id, "error", err)
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
      logger.Errorw("file data chunk expected", "error", "wrong format")
      return nil
    }

    _, err = ff.Write(chunk.Chunk)
    if err != nil {
      logger.Errorw("failed to write next file chunk in buffer", "error", err)
      return err
    }
  }

  return stream.SendAndClose(&api.SaveFileStreamResponse{
    File: &api.File{
      Id:               id,
      Name:             file.File.Name,
      CreateUtcNano:    timeCreated,
    },
  })
}