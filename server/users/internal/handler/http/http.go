package http

import (
  "net/http"
  "io"
  "time"
  "context"
  "encoding/json"

  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/users/internal/controller"
  "github.com/bd878/gallery/server/users/pkg/model"
  "github.com/bd878/gallery/server/utils"
)

type Controller interface {
  AddUser(ctx context.Context, log *logger.Logger, params *model.AddUserParams) error
  HasUser(ctx context.Context, log *logger.Logger, params *model.HasUserParams) (bool, error)
  RefreshToken(ctx context.Context, log *logger.Logger, params *model.RefreshTokenParams) error
  GetUser(ctx context.Context, log *logger.Logger, params *model.GetUserParams) (*model.User, error)
}

type Config struct {
  Domain string
}

type Handler struct {
  controller      Controller
  config          Config
}

func New(controller Controller, config Config) *Handler {
  return &Handler{controller, config}
}

func (h *Handler) Authenticate(w http.ResponseWriter, req *http.Request) {
  userName, ok := getName(w, req)
  if !ok {
    logger.Error("cannot get name")
    return
  }

  password, ok := getPassword(w, req)
  if !ok {
    logger.Error("cannot get password")
    return
  }

  exists, err := h.controller.HasUser(context.Background(), logger.Default(), &model.HasUserParams{
    User: &model.User{
      Name: userName,
      Password: password,
    },
  })
  if err != nil {
    logger.Error("failed to check if user exists:", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if !exists {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "no user,password pair",
    }); err != nil {
      logger.Error("failed to send no user,password pair:", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    return
  }

  user, err := h.controller.GetUser(context.Background(), logger.Default(), &model.GetUserParams{
    User: &model.User{Name: userName},
  })
  switch err {
  case controller.ErrTokenExpired:
    logger.Infoln("token expired")
    err := refreshToken(h, w, userName)
    if err != nil {
      logger.Error("Cannot refresh token: ", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

  case controller.ErrNotFound:
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "no user,password pair",
    }); err != nil {
      logger.Error("cannot respond no user: ", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

  case nil:
    err := attachTokenToResponse(w, user.Token, h.config.Domain, user.ExpiresUTCNano)
    if err != nil {
      logger.Error("Cannot attach token to response: ", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

  default:
    logger.Error("Unknown error: ", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(model.ServerResponse{
    Status: "ok",
    Description: "authenticated",
  }); err != nil {
    logger.Error("cannot send ok authenticated: ", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) Auth(w http.ResponseWriter, req *http.Request) {
  cookie, err := req.Cookie("token")
  if err != nil {
    logger.Error("bad cookie")
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  logger.Infoln("cookie value", cookie.Value)
  user, err := h.controller.GetUser(context.Background(), logger.Default(), &model.GetUserParams{
    User: &model.User{Token: cookie.Value},
  })
  if err == controller.ErrTokenExpired {
    logger.Infoln("token expired")
    if err := json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
      ServerResponse: model.ServerResponse{
        Status: "ok",
        Description: "token expired",
      },
      Expired: true,
    }); err != nil {
      logger.Error("cannot send token expired error: ", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
  }
  if err != nil {
    logger.Error("failed to get user by token: ", err)
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
      ID: -1,
      Name: user.Name,
      Token: user.Token,
      ExpiresUTCNano: user.ExpiresUTCNano,
    },
  }); err != nil {
    logger.Error("failed to send authorize response: ", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) Register(w http.ResponseWriter, req *http.Request) {
  userName, ok := getName(w, req)
  if !ok {
    logger.Error("cannot get user name")
    return
  }

  exists, err := h.controller.HasUser(context.Background(), logger.Default(), &model.HasUserParams{
    User: &model.User{Name: userName},
  })
  if err != nil {
    logger.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if exists {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "name exists",
    }); err != nil {
      logger.Error(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    return
  }

  password, ok := getPassword(w, req)
  if !ok {
    logger.Error("cannot get password from request")
    return
  }

  token, expiresUtcNano := createToken(w, h.config.Domain)

  logger.Infoln("user, password, token, expires", userName, password, token, expiresUtcNano)
  err = h.controller.AddUser(context.Background(), logger.Default(), &model.AddUserParams{
    User: &model.User{
      Name:                  userName,
      Password:              password,
      Token:                 token,
      ExpiresUTCNano:        expiresUtcNano,
    },
  })
  if err != nil {
    logger.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if err := json.NewEncoder(w).Encode(model.ServerResponse{
    Status: "ok",
    Description: "created",
  }); err != nil {
    logger.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func (h *Handler) ReportStatus(w http.ResponseWriter, _ *http.Request) {
  if _, err := io.WriteString(w, "ok\n"); err != nil {
    logger.Error(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func getName(w http.ResponseWriter, req *http.Request) (name string, ok bool) {
  return getTextField(w, req, "name")
}

func getPassword(w http.ResponseWriter, req *http.Request) (pass string, ok bool) {
  return getTextField(w, req, "password")
}

func getTextField(w http.ResponseWriter, req *http.Request, field string) (value string, ok bool) {
  value = req.PostFormValue(field)
  if value == "" {
    if err := json.NewEncoder(w).Encode(model.ServerResponse{
      Status: "ok",
      Description: "no " + field,
    }); err != nil {
      logger.Error(err)
      w.WriteHeader(http.StatusInternalServerError)
    }
    ok = false
  } else {
    ok = true
  }
  return
}

func refreshToken(h *Handler, w http.ResponseWriter, userName string) error {
  token, expiresUtcNano := createToken(w, h.config.Domain)

  return h.controller.RefreshToken(context.Background(), logger.Default(), &model.RefreshTokenParams{
    User: &model.User{
      Name:               userName,
      Token:              token,
      ExpiresUTCNano:     expiresUtcNano,
    },
  })
}

func createToken(w http.ResponseWriter, domain string) (token string, expires int64) {
  token = utils.RandomString(10)
  expiresAt := time.Now().Add(time.Hour * 24 * 5)

  http.SetCookie(w, &http.Cookie{
    Name:             "token",
    Value:             token,
    Domain:            domain,
    Expires:           expiresAt,
    Path:             "/",
    HttpOnly:          true,
  })

  return token, expiresAt.UnixNano()
}

func attachTokenToResponse(w http.ResponseWriter, token, domain string, expiresUtcNano int64) (err error) {
  http.SetCookie(w, &http.Cookie{
    Name:          "token",
    Value:          token,
    Domain:         domain,
    Expires:        time.Unix(0, expiresUtcNano),
    Path:          "/",
    HttpOnly:       true,
  })

  return
}

