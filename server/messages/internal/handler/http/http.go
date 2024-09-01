package http

import (
  "log"
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
)

const selectNoLimit int = -1

type userGateway interface {
  Auth(ctx context.Context, token string) (*usermodel.User, error)
}

type Controller interface {
  SaveMessage(ctx context.Context, msg *model.Message) (*model.Message, error)
  ReadUserMessages(
    ctx context.Context,
    userId usermodel.UserId,
    limit int32,
    offset int32,
    asc bool,
  ) (
    *model.MessagesList,
    error,
  )
}

type Handler struct {
  ctrl Controller
  userGateway userGateway
  dataPath string
}

func New(
  ctrl Controller,
  userGateway userGateway,
  dataPath string,
) *Handler {
  return &Handler{ctrl, userGateway, dataPath}
}

func (h *Handler) CheckAuth(
  next func (w http.ResponseWriter, req *http.Request),
) func (w http.ResponseWriter, req *http.Request) {
  return func(w http.ResponseWriter, req *http.Request) {
    cookie, err := req.Cookie("token")
    if err != nil {
      log.Println("bad cookie")
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    log.Println("cookie value =", cookie.Value)
    user, err := h.userGateway.Auth(context.Background(), cookie.Value)
    if err != nil {
      log.Println(err) // TODO: return invalid token response instead
      if err := json.NewEncoder(w).Encode(model.ServerResponse{
        Status: "ok",
        Description: "token not found",
      }); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
      }
      return
    }

    log.Println("request for user id, name, token =", user.Id, user.Name, user.Token)

    req = req.WithContext(
      context.WithValue(context.Background(), userContextKey{}, user),
    )

    next(w, req)
  }
}

type userContextKey struct {}

func (h *Handler) SendMessage(w http.ResponseWriter, req *http.Request) {
  var err error
  if err = req.ParseMultipartForm(1); err != nil {
    log.Println(err)
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
      log.Println(err)
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
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    if _, err := io.Copy(ff, f); err != nil {
      log.Println(err)
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
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
    }
    return
  }

  var msg *model.Message
  if msg, err = h.ctrl.SaveMessage(context.Background(), &model.Message{
    UserId: user.Id,
    Value: value,
    FileName: fileName,
    FileId: model.FileId(fileId),
  }); err != nil {
    log.Println(err)
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
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) ReadMessages(w http.ResponseWriter, req *http.Request) {
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
        log.Println(err)
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
          log.Println(err)
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
        log.Println(err)
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

  res, err := h.ctrl.ReadUserMessages(
    context.Background(),
    usermodel.UserId(user.Id),
    int32(limitInt),
    int32(offsetInt),
    ascending,
  )
  if err != nil {
    log.Println(err)
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
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) ReadFile(w http.ResponseWriter, req *http.Request) {
  values := req.URL.Query()
  filename := values.Get("id")
  if filename == "" {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "empty file id",
    }); err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
    }
    return
  }

  ff, err := os.Open(filepath.Join(h.dataPath, filename))
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  mimetype := mime.TypeByExtension(filepath.Ext(filename))
  if mimetype == "" {
    mimetype = "text/plain"
  }

  w.Header().Set("Content-Type", mimetype)

  if _, err := io.Copy(w, ff); err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) GetStatus(w http.ResponseWriter, _ *http.Request) {
  if _, err := io.WriteString(w, "ok\n"); err != nil {
    log.Println(err)
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
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
    }
    return nil, false
  }
  return user, true
}