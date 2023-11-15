package fcgi

import (
  "log"
  "net/http"
  "io"
  "time"
  "context"
  "encoding/json"

  "github.com/bd878/gallery/server/user/internal/controller/users"
  "github.com/bd878/gallery/server/user/pkg/model"
  "github.com/bd878/gallery/server/utils"
)

type Handler struct {
  ctrl *users.Controller
}

func New(ctrl *users.Controller) *Handler {
  return &Handler{ctrl}
}

func (h *Handler) AuthUser(w http.ResponseWriter, req *http.Request) {
  userName := req.PostFormValue("name")
  if userName == "" {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "keine name",
    }); err != nil {
      log.Println(err)
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
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
    }
    return
  }

  token := utils.RandomString(10)
  http.SetCookie(w, &http.Cookie{
    Name: "token",
    Value: token,
    Path: "/",
    Domain: "127.0.0.1:8080/",
    Expires: time.Now().Add(time.Hour * 24 * 5),
    HttpOnly: true,
  })

  log.Println("user, password, token=", userName, password, token)
  if err := h.ctrl.Add(context.Background(), &model.User{
    Name: userName,
    Password: password,
    Token: token,
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

func (h *Handler) ReportStatus(w http.ResponseWriter, _ *http.Request) {
  log.Println("report ok status")
  if _, err := io.WriteString(w, "ok\n"); err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}