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

  usermodel "github.com/bd878/gallery/server/users/pkg/model"
  "github.com/bd878/gallery/server/messages/pkg/model"
  "github.com/bd878/gallery/server/utils"
  "github.com/bd878/gallery/server/logger"
)

const selectNoLimit int = -1

type userGateway interface {
  Auth(ctx context.Context, log *logger.Logger, params *model.AuthParams) (*usermodel.User, error)
}

type Controller interface {
  SaveMessage(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) (*model.Message, error)
  ReadUserMessages(ctx context.Context, log *logger.Logger, params *model.ReadUserMessagesParams) (*model.MessagesList, error)
}

type Handler struct {
  controller Controller
  userGateway userGateway
  dataPath string
}

func New(controller Controller, userGateway userGateway, dataPath string) *Handler {
  return &Handler{controller, userGateway, dataPath}
}

func (h *Handler) CheckAuth(next func (ctx context.Context, log *logger.Logger, w http.ResponseWriter, req *http.Request)) (
  func (ctx context.Context, log *logger.Logger, w http.ResponseWriter, req *http.Request),
) {
  return func(ctx context.Context, log *logger.Logger, w http.ResponseWriter, req *http.Request) {
    cookie, err := req.Cookie("token")
    if err != nil {
      log.Errorln("bad cookie")
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    log.Infoln("cookie value", cookie.Value)
    user, err := h.userGateway.Auth(context.Background(), log, &model.AuthParams{Token: cookie.Value})
    if err != nil {
      logger.Errorln(err) // TODO: return invalid token response instead
      if err := json.NewEncoder(w).Encode(model.ServerResponse{
        Status: "ok",
        Description: "token not found",
      }); err != nil {
        logger.Error(err)
        w.WriteHeader(http.StatusInternalServerError)
      }
      return
    }

    log.Infoln("user id", user.Id, "name", user.Name, "token", user.Token)

    req = req.WithContext(
      context.WithValue(context.Background(), userContextKey{}, user),
    )

    next(ctx, log, w, req)
  }
}

type userContextKey struct {}

func (h *Handler) SendMessage(ctx context.Context, log *logger.Logger, w http.ResponseWriter, req *http.Request) {
  var err error
  if err = req.ParseMultipartForm(1); err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  var user *usermodel.User
  var ok bool
  if user, ok = getUser(w, req); !ok {
    return
  }

  var fileName string
  var fileId string
  if _, ok := req.MultipartForm.File["file"]; ok {
    f, fh, err := req.FormFile("file")
    if err != nil {
      log.Error(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    fileName = filepath.Base(fh.Filename)
    fileId = strings.ToLower(utils.RandomString(10) + filepath.Ext(fh.Filename))

    ff, err := os.OpenFile(
      filepath.Join(h.dataPath, fileId),
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

  value := req.PostFormValue("message")

  if fileName == "" && value == "" {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "empty fields",
    }); err != nil {
      log.Error(err)
      w.WriteHeader(http.StatusInternalServerError)
    }
    return
  }

  var msg *model.Message
  if msg, err = h.controller.SaveMessage(context.Background(), log, &model.SaveMessageParams{
    Message: &model.Message{
      UserId: int(user.Id),
      Value: value,
      FileName: fileName,
      FileId: model.FileId(fileId),
    },
  }); err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(model.NewMessageServerResponse{
    ServerResponse: model.ServerResponse{
      Status: "ok",
      Description: "accepted",
    },
    Message: *msg,
  }); err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) ReadMessages(ctx context.Context, log *logger.Logger, w http.ResponseWriter, req *http.Request) {
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
      if err = json.NewEncoder(w).Encode(model.ServerResponse{
        Status: "ok",
        Description: fmt.Sprintf("wrong \"%s\" query param", "limit"),
      }); err != nil {
        log.Error(err)
        w.WriteHeader(http.StatusInternalServerError)
      }
      return
    }

    if values.Has("offset") {
      offsetInt, err = strconv.Atoi(values.Get("offset"))
      if err != nil {
        if err = json.NewEncoder(w).Encode(model.ServerResponse{
          Status: "ok",
          Description: fmt.Sprintf("wrong \"%s\" query param", "offset"),
        }); err != nil {
          log.Error(err)
          w.WriteHeader(http.StatusInternalServerError)
        }
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
      if err = json.NewEncoder(w).Encode(model.ServerResponse{
        Status: "ok",
        Description: fmt.Sprintf("wrong \"%s\" query param", "asc"),
      }); err != nil {
        log.Error(err)
        w.WriteHeader(http.StatusInternalServerError)
      }
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
      UserId: usermodel.UserId(user.Id),
      Limit: int32(limitInt),
      Offset: int32(offsetInt),
      Ascending: ascending,
    },
  )
  if err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(model.MessagesListServerResponse{
    ServerResponse: model.ServerResponse{
      Status: "ok",
    },
    Messages: res.Messages,
    IsLastPage: res.IsLastPage,
  }); err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) ReadFile(ctx context.Context, log *logger.Logger, w http.ResponseWriter, req *http.Request) {
  values := req.URL.Query()
  filename := values.Get("id")
  if filename == "" {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "empty file id",
    }); err != nil {
      log.Error(err)
      w.WriteHeader(http.StatusInternalServerError)
    }
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
    mimetype = "text/plain"
  }

  w.Header().Set("Content-Type", mimetype)

  if _, err := io.Copy(w, ff); err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) GetStatus(ctx context.Context, log *logger.Logger, w http.ResponseWriter, _ *http.Request) {
  if _, err := io.WriteString(w, "ok\n"); err != nil {
    log.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func getUser(w http.ResponseWriter, req *http.Request) (*usermodel.User, bool) {
  user, ok := req.Context().Value(userContextKey{}).(*usermodel.User)
  if !ok {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "user required",
    }); err != nil {
      logger.Error(err)
      w.WriteHeader(http.StatusInternalServerError)
    }
    return nil, false
  }
  return user, true
}