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

/* TODO: rewrite global config on singletone pattern */
type Config struct {
  Domainname string
}

type Handler struct {
  ctrl *users.Controller
  cfg   Config
}

func New(ctrl *users.Controller, cfg Config) *Handler {
  return &Handler{ctrl, cfg}
}

func (h *Handler) Authenticate(w http.ResponseWriter, req *http.Request) {
  var userName, password string 
  var ok, exists bool
  var user *model.User
  var err error

  if userName, ok = getName(w, req); !ok {
    log.Println("cannot get name")
    return
  }

  if password, ok = getPassword(w, req); !ok {
    log.Println("cannot get password")
    return
  }

  exists, err = h.ctrl.Has(context.Background(), &model.User{Name: userName, Password: password})
  if err != nil {
    log.Println("failed to check if user exists:", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if !exists {
    if err = json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "no user,password pair",
    }); err != nil {
      log.Println("failed to send no user,password pair:", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    return
  }

  user, err = h.ctrl.Get(context.Background(), &model.User{Name: userName})
  switch err {
  case controller.ErrTokenExpired:
    log.Println("token expired")
    err = refreshToken(h, w, userName)
    if err != nil {
      log.Println("Cannot refresh token: ", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

  case controller.ErrNotFound:
    if err = json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "no user,password pair",
    }); err != nil {
      log.Println("cannot respond no user: ", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

  case nil:
    err = attachTokenToResponse(w, user.Token, h.cfg.Domainname, user.Expires)
    if err != nil {
      log.Println("Cannot attach token to response: ", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

  default:
    log.Println("Unknown error: ", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err = json.NewEncoder(w).Encode(model.ServerResponse{
    Status: "ok",
    Description: "authenticated",
  }); err != nil {
    log.Println("cannot send ok authenticated: ", err)
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
  if err == controller.ErrTokenExpired {
    log.Println("token expired")
    if err := json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
      ServerResponse: model.ServerResponse{
        Status: "ok",
        Description: "token expired",
      },
      Expired: true,
    }); err != nil {
      log.Println("cannot send token expired error: ", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
  }
  // TODO: check for token expired error
  if err != nil {
    log.Println("failed to get user by token: ", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
    ServerResponse: model.ServerResponse{
      Status: "ok",
      Description: "token valid",
    },
    Expired: false,
    User: model.User{
      Id: -1,
      Name: user.Name,
      Token: user.Token,
      Expires: user.Expires,
    },
  }); err != nil {
    log.Println("failed to send authorize response: ", err)
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

  token, expires := createToken(w, h.cfg.Domainname)

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

func refreshToken(h *Handler, w http.ResponseWriter, userName string) (err error) {
  var token, expires string

  token, expires = createToken(w, h.cfg.Domainname)

  return h.ctrl.Refresh(context.Background(), &model.User{Name: userName, Token: token, Expires: expires})
}

func createToken(w http.ResponseWriter, domain string) (token string, expires string) {
  token = utils.RandomString(10)
  expiresAt := time.Now().Add(time.Hour * 24 * 5)
  expiresRaw, err := expiresAt.MarshalText()
  if err != nil {
    log.Println("failed to marshal text: ", err)
  }

  expires = string(expiresRaw)

  http.SetCookie(w, &http.Cookie{
    Name: "token",
    Value: token,
    Domain: domain,
    Expires: expiresAt,
    Path: "/",
    HttpOnly: true,
  })

  return
}

func attachTokenToResponse(w http.ResponseWriter, token, domain string, expires string) (err error) {
  var tokenExpiresTime time.Time

  err = tokenExpiresTime.UnmarshalText([]byte(expires))
  if err != nil {
    return
  }

  http.SetCookie(w, &http.Cookie{
    Name: "token",
    Value: token,
    Domain: domain,
    Expires: tokenExpiresTime,
    Path: "/",
    HttpOnly: true,
  })

  return
}

