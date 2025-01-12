package middleware

import (
  "context"
  "net/http"
  "encoding/json"

  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/pkg/model"
  usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

func Log(next Handler) Handler {
  return func(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
    next(log, w, req)
  }
}

type UserContextKey struct {}

type userGateway interface {
  Auth(ctx context.Context, log *logger.Logger, params *model.AuthParams) (*usermodel.User, error)
}

type AuthBuilder struct {
  Gateway userGateway
}

func (b *AuthBuilder) Auth(next Handler) Handler {
  return func(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
    cookie, err := req.Cookie("token")
    if err != nil {
      log.Errorln("bad cookie")
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    log.Infoln("cookie value", cookie.Value)
    user, err := b.Gateway.Auth(context.Background(), log, &model.AuthParams{Token: cookie.Value})
    if err != nil {
      logger.Errorln(err) // TODO: return invalid token response instead
      if err := json.NewEncoder(w).Encode(model.ServerResponse{
        Status: "ok",
        Description: "token not found",
      }); err != nil {
        logger.Error(err)
        w.WriteHeader(http.StatusInternalServerError)
      }
      return
    }

    log.Infoln("user id", user.ID, "name", user.Name, "token", user.Token)

    req = req.WithContext(
      context.WithValue(context.Background(), UserContextKey{}, user),
    )

    next(log, w, req)
  }
}
