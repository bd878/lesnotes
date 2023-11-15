package fcgi

import (
  "log"
  "net/http"
  "os"
  "io"
  "context"
  "path/filepath"
  "encoding/json"

  "github.com/bd878/gallery/server/messages/internal/controller/messages"
  "github.com/bd878/gallery/server/messages/pkg/model"
)

type Handler struct {
  ctrl *messages.Controller
  dataPath string
}

func New(
  ctrl *messages.Controller,
  dataPath string,
) *Handler {
  return &Handler{ctrl, dataPath}
}

func (h *Handler) SaveMessage(w http.ResponseWriter, req *http.Request) {
  if err := req.ParseMultipartForm(1); err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
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

    ff, err := os.OpenFile(
      filepath.Join(h.dataPath, fh.Filename),
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
    filename = fh.Filename
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
  v, err := h.ctrl.ReadAllMessages(context.Background())
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  log.Println("send n'th messages =", len(v))
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