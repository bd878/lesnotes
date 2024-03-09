package fcgi

import (
  "log"
  "net/http"
  "os"
  "io"
  "context"
  "strings"
  "path/filepath"
  "encoding/json"

  "github.com/bd878/gallery/server/utils"
  usermodel "github.com/bd878/gallery/server/user/pkg/model"
  "github.com/bd878/gallery/server/messages/pkg/model"
)

type userGateway interface {
  Auth(ctx context.Context, token string) (*usermodel.User, error)
}

type Controller interface {
  SaveMessage(ctx context.Context, msg *model.Message) error
  ReadUserMessages(ctx context.Context, userId usermodel.UserId) ([]model.Message, error)
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

func (h *Handler) SaveMessage(w http.ResponseWriter, req *http.Request) {
  if err := req.ParseMultipartForm(1); err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  var user *usermodel.User
  var ok bool
  if user, ok = getUser(w, req); !ok {
    return
  }

  var filename string
  if _, ok := req.MultipartForm.File["file"]; ok {
    f, fh, err := req.FormFile("file")
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    filename = strings.ToLower(utils.RandomString(10) + filepath.Ext(fh.Filename))
    log.Println("filename=", filename)

    ff, err := os.OpenFile(
      filepath.Join(h.dataPath, filename),
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

  msg := req.PostFormValue("message")

  if filename == "" && msg == "" {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "empty fields",
    }); err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
    }
    return
  }

  if err := h.ctrl.SaveMessage(context.Background(), &model.Message{
    UserId: user.Id,
    Value: msg,
    File: filename,
  }); err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(model.ServerResponse{
    Status: "ok",
    Description: "accepted",
  }); err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) ReadMessages(w http.ResponseWriter, req *http.Request) {
  var user *usermodel.User
  var ok bool
  if user, ok = getUser(w, req); !ok {
    return
  }

  v, err := h.ctrl.ReadUserMessages(context.Background(), usermodel.UserId(user.Id))
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(v); err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) ReportStatus(w http.ResponseWriter, _ *http.Request) {
  log.Println("report ok status")
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