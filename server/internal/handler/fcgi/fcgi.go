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
  "github.com/bd878/gallery/server/internal/controller/users"
  "github.com/bd878/gallery/server/pkg/model"
)

type Handler struct {
  msgCtrl *messages.Controller
  usrCtrl *users.Controller
  dataPath string
}

func New(
  msgCtrl *messages.Controller,
  usrCtrl *users.Controller,
  dataPath string,
) *Handler {
  return &Handler{msgCtrl, usrCtrl, dataPath}
}

func (h *Handler) AuthUser(w http.ResponseWriter, req *http.Request) {
  userName := req.PostFormValue("name")
  if userName == "" {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "keine name",
    }); err != nil {
      w.WriteHeader(http.StatusInternalServerError)
    }
    return
  }

  password := req.PostFormValue("password")
  if password == "" {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "keine password",
    }); err != nil {
      w.WriteHeader(http.StatusInternalServerError)
    }
    return
  }

  log.Println("user, password=", userName, password)
  if err := h.usrCtrl.Add(context.Background(), &model.User{
    Name: userName,
    Password: password,
  }); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(model.ServerResponse{
    Status: "ok",
    Description: "accepted",
  }); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) SaveMessage(w http.ResponseWriter, req *http.Request) {
  var filename string
  if _, ok := req.MultipartForm.File["file"]; ok {
    f, fh, err := req.FormFile("file")
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    ff, err := os.OpenFile(
      filepath.Join(h.dataPath, fh.Filename),
      os.O_WRONLY|os.O_CREATE, 0666,
    )
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    if _, err := io.Copy(ff, f); err != nil {
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
      w.WriteHeader(http.StatusInternalServerError)
    }
    return
  }

  if err := h.msgCtrl.SaveMessage(context.Background(), &model.Message{
    Value: msg,
    File: filename,
  }); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(model.ServerResponse{
    Status: "ok",
    Description: "accepted",
  }); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) ReadMessages(w http.ResponseWriter, req *http.Request) {
  v, err := h.msgCtrl.ReadAllMessages(context.Background())
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