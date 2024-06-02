package http

import (
  "log"
  "net/http"
  "io"
  "time"
  "context"
  "encoding/json"

  "github.com/bd878/gallery/server/users/internal/controller"
  "github.com/bd878/gallery/server/users/internal/controller/users"
  "github.com/bd878/gallery/server/users/pkg/model"
  "github.com/bd878/gallery/server/utils"
)

type Handler struct {
  ctrl *users.Controller
}

func New(ctrl *users.Controller) *Handler {
  return &Handler{ctrl}
}

func (h *Handler) Authenticate(w http.ResponseWriter, req *http.Request) {
  var userName, password string 
  var ok bool

  if userName, ok = getName(w, req); !ok {
    return
  }

  if password, ok = getPassword(w, req); !ok {
    return
  }

  exists, err := h.ctrl.Has(context.Background(), &model.User{Name: userName, Password: password})
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if !exists {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "no user,password pair",
    }); err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    return
  }

  token, expires := createToken(w)
  err = h.ctrl.Refresh(context.Background(), &model.User{Name: userName, Token: token, Expires: expires})
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(model.ServerResponse{
    Status: "ok",
    Description: "authenticated",
  }); err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) Auth(w http.ResponseWriter, req *http.Request) {
  cookie, err := req.Cookie("token")
  if err != nil {
    log.Println("bad cookie")
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  log.Println("cookie value =", cookie.Value)
  user, err := h.ctrl.Get(context.Background(), &model.User{Token: cookie.Value})
  if err == controller.ErrTokenInvalid {
    if err := json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
      ServerResponse: model.ServerResponse{
        Status: "ok",
        Description: "invalid token",
      },
      Valid: false,
    }); err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
  }
  // TODO: check for token expired error
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
    ServerResponse: model.ServerResponse{
      Status: "ok",
      Description: "token valid",
    },
    Valid: true,
    User: model.User{
      Id: -1,
      Name: user.Name,
      Token: user.Token,
      Expires: user.Expires,
    },
  }); err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) Register(w http.ResponseWriter, req *http.Request) {
  var userName, password string
  var ok bool

  // user name
  if userName, ok = getName(w, req); !ok {
    return
  }

  // already exists
  exists, err := h.ctrl.Has(context.Background(), &model.User{Name: userName})
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if exists {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "name exists",
    }); err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    return
  }

  // password
  if password, ok = getPassword(w, req); !ok {
    return
  }

  token, expires := createToken(w)

  log.Println("user, password, token, expires=", userName, password, token, expires)
  if err := h.ctrl.Add(context.Background(), &model.User{
    Name: userName,
    Password: password,
    Token: token,
    Expires: expires,
  }); err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(model.ServerResponse{
    Status: "ok",
    Description: "created",
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

func getName(w http.ResponseWriter, req *http.Request) (name string, ok bool) {
  name, ok = getTextField(w, req, "name")
  return
}

func getPassword(w http.ResponseWriter, req *http.Request) (pass string, ok bool) {
  pass, ok = getTextField(w, req, "password")
  return
}

func getTextField(w http.ResponseWriter, req *http.Request, field string) (value string, ok bool) {
  value = req.PostFormValue(field)
  if value == "" {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "no " + field,
    }); err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
    }
    ok = false
  } else {
    ok = true
  }
  return
}

func createToken(w http.ResponseWriter) (token string, expires string) {
  token = utils.RandomString(10)
  expiresAt := time.Now().Add(time.Hour * 24 * 5)
  expires = expiresAt.String()

  http.SetCookie(w, &http.Cookie{
    Name: "token",
    Value: token,
    Domain: "galleryexample.com",
    Expires: expiresAt,
    Path: "/",
    HttpOnly: true,
  })
  return
}