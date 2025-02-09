package http

import (
  "net/http"
  "strconv"
  "io"
  "fmt"
  "context"
  "path/filepath"
  "encoding/json"

  filesmodel "github.com/bd878/gallery/server/files/pkg/model"
  servermodel "github.com/bd878/gallery/server/pkg/model"
  "github.com/bd878/gallery/server/messages/pkg/model"
  "github.com/bd878/gallery/server/utils"
  "github.com/bd878/gallery/server/logger"
)

const selectNoLimit int = -1

type Controller interface {
  SaveMessage(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) (*model.SaveMessageResult, error)
  UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) (*model.UpdateMessageResult, error)
  DeleteMessage(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) (*model.DeleteMessageResult, error)
  ReadUserMessages(ctx context.Context, log *logger.Logger, params *model.ReadUserMessagesParams) (*model.ReadUserMessagesResult, error)
}

// TODO: move files into controller?
type FilesGateway interface {
  SaveFileStream(ctx context.Context, log *logger.Logger, stream io.Reader, params *model.SaveFileParams) (*model.SaveFileResult, error)
  ReadBatchFiles(ctx context.Context, log *logger.Logger, params *model.ReadBatchFilesParams) (*model.ReadBatchFilesResult, error)
  // TODO: DeleteFile
}

type Handler struct {
  controller    Controller
  filesGateway  FilesGateway
}

func New(controller Controller, filesGateway FilesGateway) *Handler {
  return &Handler{controller, filesGateway}
}

func (h *Handler) SendMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
  var err error
  if err = req.ParseMultipartForm(1); err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  user, ok := utils.GetUser(w, req)
  if !ok {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(servermodel.ServerResponse{
      Status: "ok",
      Description: "user required",
    })
    return
  }

  var fileName string
  var fileID int32
  if _, ok := req.MultipartForm.File["file"]; ok {
    f, fh, err := req.FormFile("file")
    if err != nil {
      log.Error(err)
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(servermodel.ServerResponse{
        Status: "error",
        Description: "cannot read file",
      })
      return
    }

    fileName = filepath.Base(fh.Filename)

    // TODO: create fileID here, pipe in goroutine
    fileResult, err := h.filesGateway.SaveFileStream(context.Background(), log, f, &model.SaveFileParams{
      UserID: user.ID,
      Name:   fileName,
    })
    if err != nil {
      log.Error(err)
      w.WriteHeader(http.StatusInternalServerError)
      json.NewEncoder(w).Encode(servermodel.ServerResponse{
        Status: "error",
        Description: "cannot save file",
      })
      return
    }

    fileID = fileResult.ID
  }

  text := req.PostFormValue("text")

  log.Infoln("text=", text, "user name=", user.Name)

  if fileName == "" && text == "" {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(servermodel.ServerResponse{
      Status: "error",
      Description: "text or file are empty",
    })
    return
  }

  resp, err := h.controller.SaveMessage(req.Context(), log, &model.SaveMessageParams{
    Message: &model.Message{
      UserID: user.ID,
      Text:   text,
      File:   &filesmodel.File{
        ID:     fileID,
      },
    },
  })
  if err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  json.NewEncoder(w).Encode(model.NewMessageServerResponse{
    ServerResponse: servermodel.ServerResponse{
      Status: "ok",
      Description: "accepted",
    },
    Message: model.Message{
      ID:            resp.ID,
      File:          &filesmodel.File{
        ID:          fileID,
        Name:        fileName,
      },
      Text:          text,
      UpdateUTCNano: resp.UpdateUTCNano,
      CreateUTCNano: resp.CreateUTCNano,
    },
  })
}

func (h *Handler) DeleteMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
  user, ok := utils.GetUser(w, req)
  if !ok {
    log.Error("user not found")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(servermodel.ServerResponse{
      Status: "ok",
      Description: "user required",
    })
    return
  }

  values := req.URL.Query()
  if values.Get("id") == "" {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(servermodel.ServerResponse{
      Status: "ok",
      Description: "empty message id",
    })
    return
  }

  id, err := strconv.Atoi(values.Get("id"))
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(servermodel.ServerResponse{
      Status: "ok",
      Description: "invalid id",
    })
    return
  }

  _, err = h.controller.DeleteMessage(req.Context(), log, &model.DeleteMessageParams{
    ID: int32(id),
    UserID: user.ID,
  })
  if err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  json.NewEncoder(w).Encode(model.DeleteMessageServerResponse{
    ServerResponse: servermodel.ServerResponse{
      Status: "ok",
      Description: "deleted",
    },
    ID: int32(id),
  })
}

func (h *Handler) UpdateMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
  user, ok := utils.GetUser(w, req)
  if !ok {
    log.Error("user not found")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(servermodel.ServerResponse{
      Status: "ok",
      Description: "user required",
    })
    return
  }

  values := req.URL.Query()
  if values.Get("id") == "" {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(servermodel.ServerResponse{
      Status: "ok",
      Description: "empty message id",
    })
    return
  }

  id, err := strconv.Atoi(values.Get("id"))
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(servermodel.ServerResponse{
      Status: "ok",
      Description: "invalid id",
    })
    return
  }

  text := req.PostFormValue("text")
  if text == "" {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(servermodel.ServerResponse{
      Status: "ok",
      Description: "empty message field",
    })
    return
  }

  resp, err := h.controller.UpdateMessage(req.Context(), log, &model.UpdateMessageParams{
    ID:     int32(id),
    UserID: user.ID,
    Text:   text,
  })
  if err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  json.NewEncoder(w).Encode(model.UpdateMessageServerResponse{
    ServerResponse: servermodel.ServerResponse{
      Status: "ok",
      Description: "accepted",
    },
    ID: resp.ID,
    UpdateUTCNano: resp.UpdateUTCNano,
  })
}

func (h *Handler) ReadMessages(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
  var limitInt, offsetInt, orderInt int
  var ascending bool
  var err error

  user, ok := utils.GetUser(w, req)
  if !ok {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(servermodel.ServerResponse{
      Status: "ok",
      Description: "user required",
    })
    return
  }

  values := req.URL.Query()

  if values.Has("limit") {
    limitInt, err = strconv.Atoi(values.Get("limit"))
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(servermodel.ServerResponse{
        Status: "error",
        Description: fmt.Sprintf("wrong \"%s\" query param", "limit"),
      })
      return
    }

    if values.Has("offset") {
      offsetInt, err = strconv.Atoi(values.Get("offset"))
      if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(servermodel.ServerResponse{
          Status: "error",
          Description: fmt.Sprintf("wrong \"%s\" query param", "offset"),
        })
        return
      }
    } else {
      offsetInt = 0
    }
  } else {
    limitInt = selectNoLimit
  }

  if values.Has("asc") {
    orderInt, err = strconv.Atoi(values.Get("asc"))
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(servermodel.ServerResponse{
        Status: "error",
        Description: fmt.Sprintf("wrong \"%s\" query param", "asc"),
      })
      return
    }

    switch orderInt {
    case 0:
      ascending = false
    case 1:
      ascending = true
    default:
      ascending = true
    }
  } else {
    ascending = true
  }

  res, err := h.controller.ReadUserMessages(context.Background(), log, &model.ReadUserMessagesParams{
    UserID:    user.ID,
    Limit:     int32(limitInt),
    Offset:    int32(offsetInt),
    Ascending: ascending,
  })
  if err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  fileIds := make([]int32, len(res.Messages))
  for i, message := range res.Messages {
    fileIds[i] = message.File.ID
  }

  filesRes, err := h.filesGateway.ReadBatchFiles(context.Background(), log, &model.ReadBatchFilesParams{
    UserID: user.ID,
    IDs:    fileIds,
  })
  if err != nil {
    log.Error("failed to read batch files", "error=", err)
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(servermodel.ServerResponse{
      Status: "error",
      Description: "failed to read files",
    })
    return
  }

  for _, message := range res.Messages {
    message.File.Name = filesRes.Files[message.File.ID].Name
  }

  json.NewEncoder(w).Encode(model.MessagesListServerResponse{
    ServerResponse: servermodel.ServerResponse{
      Status: "ok",
    },
    Messages: res.Messages,
    IsLastPage: res.IsLastPage,
  })
}

func (h *Handler) GetStatus(log *logger.Logger, w http.ResponseWriter, _ *http.Request) {
  if _, err := io.WriteString(w, "ok\n"); err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}
