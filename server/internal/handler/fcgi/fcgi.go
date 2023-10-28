package fcgi

import (
  "net/http"
  "io"
  "context"
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
  err := req.ParseForm()
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  msg := req.PostFormValue("message")
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

  if err := json.NewEncoder(w).Encode(v); err != nil {
    panic(err)
  }
}

func (h *Handler) ReportStatus(w http.ResponseWriter, _ *http.Request) {
  if _, err := io.WriteString(w, "ok"); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
  }
}