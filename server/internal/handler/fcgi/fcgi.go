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
    log.Println("err =", err)
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  log.Println("filename, size:", fh.Filename, fh.Size)
  switch f.(type) {
  case *os.File:
    log.Println("file is of type *os.File")
    info, err := f.(*os.File).Stat()
    if err != nil {
      panic(err)
    }
    log.Println(filepath.Join(os.TempDir(), info.Name()))
  default:
    log.Println("file is of some other type")
  }

  err = h.ctrl.SaveMessage(context.Background(), msg)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
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