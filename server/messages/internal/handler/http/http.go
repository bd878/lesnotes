package http

import (
  "net/http"
  "strconv"
  "strings"
  "os"
  "io"
  "fmt"
  "context"
  "mime"
  "path/filepath"
  "encoding/json"

  httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
  usermodel "github.com/bd878/gallery/server/users/pkg/model"
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

type Handler struct {
  controller Controller
  dataPath string
}

func New(controller Controller, dataPath string) *Handler {
  return &Handler{controller, dataPath}
}

func (h *Handler) SendMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
  var err error
  if err = req.ParseMultipartForm(1); err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  user, ok := getUser(w, req)
  if !ok {
    return
  }

  var fileName string
  var fileID int32
  var fileUID string
  if _, ok := req.MultipartForm.File["file"]; ok {
    f, fh, err := req.FormFile("file")
    if err != nil {
      log.Error(err)
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(model.ServerResponse{
        Status: "error",
        Description: "cannot read file",
      })
      return
    }

    fileName = filepath.Base(fh.Filename)
    fileID = utils.RandomID()
    fileUID = strings.ToLower(utils.RandomString(10) + filepath.Ext(fh.Filename))

    ff, err := os.OpenFile(
      filepath.Join(h.dataPath, fileUID),
      os.O_WRONLY|os.O_CREATE, 0666,
    )
    if err != nil {
      log.Error(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    if _, err := io.Copy(ff, f); err != nil {
      log.Error(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
  }

  text := req.PostFormValue("text")

  log.Infoln("text=", text, "user name=", user.Name)

  if fileName == "" && text == "" {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "error",
      Description: "text or file are empty",
    })
    return
  }

  resp, err := h.controller.SaveMessage(req.Context(), log, &model.SaveMessageParams{
    Message: &model.Message{
      UserID: user.ID,
      Text:   text,
      FileID: fileID,
    },
  })
  if err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  json.NewEncoder(w).Encode(model.NewMessageServerResponse{
    ServerResponse: model.ServerResponse{
      Status: "ok",
      Description: "accepted",
    },
    ID: resp.ID,
    UpdateUTCNano: resp.UpdateUTCNano,
    CreateUTCNano: resp.CreateUTCNano,
  })
}

func (h *Handler) DeleteMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
  user, ok := getUser(w, req)
  if !ok {
    log.Error("user not found")
    return
  }

  values := req.URL.Query()
  if values.Get("id") == "" {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "empty message id",
    })
    return
  }

  id, err := strconv.Atoi(values.Get("id"))
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(model.ServerResponse{
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
    ServerResponse: model.ServerResponse{
      Status: "ok",
      Description: "deleted",
    },
    ID: int32(id),
  })
}

func (h *Handler) UpdateMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
  user, ok := getUser(w, req)
  if !ok {
    log.Error("user not found")
    return
  }

  values := req.URL.Query()
  if values.Get("id") == "" {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "empty message id",
    })
    return
  }

  id, err := strconv.Atoi(values.Get("id"))
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "invalid id",
    })
    return
  }

  text := req.PostFormValue("text")
  if text == "" {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(model.ServerResponse{
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
    ServerResponse: model.ServerResponse{
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
  var ok bool
  var user *usermodel.User

  if user, ok = getUser(w, req); !ok {
    return
  }

  values := req.URL.Query()

  if values.Has("limit") {
    limitInt, err = strconv.Atoi(values.Get("limit"))
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(model.ServerResponse{
        Status: "error",
        Description: fmt.Sprintf("wrong \"%s\" query param", "limit"),
      })
      return
    }

    if values.Has("offset") {
      offsetInt, err = strconv.Atoi(values.Get("offset"))
      if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(model.ServerResponse{
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
      json.NewEncoder(w).Encode(model.ServerResponse{
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

  res, err := h.controller.ReadUserMessages(
    context.Background(),
    log,
    &model.ReadUserMessagesParams{
      UserID:    user.ID,
      Limit:     int32(limitInt),
      Offset:    int32(offsetInt),
      Ascending: ascending,
    },
  )
  if err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  json.NewEncoder(w).Encode(model.MessagesListServerResponse{
    ServerResponse: model.ServerResponse{
      Status: "ok",
    },
    Messages: res.Messages,
    IsLastPage: res.IsLastPage,
  })
}

func (h *Handler) ReadFile(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
  values := req.URL.Query()
  filename := values.Get("id")
  if filename == "" {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "empty file id",
    })
    return
  }

  ff, err := os.Open(filepath.Join(h.dataPath, filename))
  if err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  mimetype := mime.TypeByExtension(filepath.Ext(filename))
  if mimetype == "" {
    mimetype = "application/octet-stream"
  }

  w.Header().Set("Content-Type", mimetype)

  if _, err := io.Copy(w, ff); err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) GetStatus(log *logger.Logger, w http.ResponseWriter, _ *http.Request) {
  if _, err := io.WriteString(w, "ok\n"); err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func getUser(w http.ResponseWriter, req *http.Request) (*usermodel.User, bool) {
  user, ok := req.Context().Value(httpmiddleware.UserContextKey{}).(*usermodel.User)
  if !ok {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "user required",
    })
    return nil, false
  }
  return user, true
}