package fcgi

import (
  "log"
  "net/http"
  "os"
  "io"
  "context"
  "path/filepath"
  "encoding/json"

  "github.com/bd878/gallery/server/internal/controller/messages"
  "github.com/bd878/gallery/server/pkg/model"
)

type Handler struct {
  ctrl *messages.Controller
}

func New(ctrl *messages.Controller) *Handler {
  return &Handler{ctrl}
}

func (h *Handler) SaveMessage(w http.ResponseWriter, req *http.Request) {
  err := req.ParseMultipartForm(1)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  msg := req.PostFormValue("message")
  log.Println("received message =", msg)

  f, fh, err := req.FormFile("file")
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  log.Println("received file of size =", fh.Filename, fh.Size)

  var fPath string
  switch f.(type) {
  case *os.File:
    info, err := f.(*os.File).Stat()
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    fPath = filepath.Join(os.TempDir(), info.Name())
  default:
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  log.Println("file full path =", fPath)
  err = h.ctrl.SaveMessage(context.Background(), &model.Message{
    Value: msg,
    File: fPath,
  })
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) ReadMessages(w http.ResponseWriter, req *http.Request) {
  v, err := h.ctrl.ReadAllMessages(context.Background())
  if err != nil {
    panic(err)
  }

  log.Println("send n'th messages =", len(v))
  if err := json.NewEncoder(w).Encode(v); err != nil {
    panic(err)
  }
}

func (h *Handler) ReportStatus(w http.ResponseWriter, _ *http.Request) {
  log.Println("report ok status")
  if _, err := io.WriteString(w, "ok\n"); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
  }
}