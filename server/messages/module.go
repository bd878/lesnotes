package messages

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bd878/gallery/server/messages/config"
	"github.com/bd878/gallery/server/internal/system"
	"github.com/bd878/gallery/server/internal/jetstream"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/internal/ddd"

	"github.com/bd878/gallery/server/messages/internal/handler/stream"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	controller "github.com/bd878/gallery/server/messages/internal/controller/service"
	threadsgateway "github.com/bd878/gallery/server/messages/internal/gateway/threads/grpc"
	filesgateway "github.com/bd878/gallery/server/messages/internal/gateway/files/grpc"
	httphandler "github.com/bd878/gallery/server/messages/internal/handler/http"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
)

func Root(ctx context.Context, cfg config.Config, svc system.Service) (err error) {
	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)
	threadsGateway := threadsgateway.New(cfg.ThreadsServiceAddr)
	filesGateway := filesgateway.New(cfg.FilesServiceAddr)
	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(svc.Logger(), usersGateway, sessionsGateway, usermodel.PublicUserID))

	dispatcher := ddd.NewEventDispatcher[ddd.Event]()
	js := jetstream.NewStream(svc.Config().NatsStream, svc.JS(), svc.Logger())
	stream.RegisterDomainEventHandlers(dispatcher, stream.NewDomainEventHandlers(
		am.NewMessagePublisher(
			js,
		),
	))

	messagesSaved := promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_saved_count",
		},
	)

	messagesController := controller.NewMessagesController(controller.MessagesConfig{RpcAddr: cfg.MessagesServiceAddr},
		dispatcher, filesGateway, threadsGateway, messagesSaved)
	translationsController := controller.NewTranslationsController(controller.TranslationsConfig{RpcAddr: cfg.MessagesServiceAddr}, dispatcher)
	commentsController := controller.NewCommentsController(controller.CommentsConfig{RpcAddr: cfg.MessagesServiceAddr}, dispatcher)

	stream.RegisterIntegrationEventHandlers(am.NewMessageSubscriber(
			js,
		),
		stream.NewIntegrationEventHandlers(messagesController),
	)

	handler := httphandler.New(messagesController, translationsController, commentsController)

	svc.ServeMux().Handle("/messages/v1/send", middleware.Build(handler.SendMessage))
	svc.ServeMux().Handle("/messages/v1/read_path", middleware.Build(handler.ReadPath))
	svc.ServeMux().Handle("/messages/v1/read", middleware.Build(handler.ReadMessages))
	svc.ServeMux().Handle("/messages/v1/update", middleware.Build(handler.UpdateMessage))
	svc.ServeMux().Handle("/messages/v1/publish", middleware.Build(handler.PublishMessages))
	svc.ServeMux().Handle("/messages/v1/private", middleware.Build(handler.PrivateMessages))
	svc.ServeMux().Handle("/messages/v1/delete", middleware.Build(handler.DeleteMessages))

	middleware.NoAuth()
	svc.ServeMux().Handle("GET /liveness", middleware.Build(handler.GetStatus))
	svc.ServeMux().Handle("GET /metrics", promhttp.Handler())

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(svc.Logger(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	svc.ServeMux().Handle("/messages/v2/send", middleware.Build(handler.SendMessageJsonAPI))
	svc.ServeMux().Handle("/messages/v2/read", middleware.Build(handler.ReadMessagesJsonAPI))
	// TODO: /threads/v2/read
	svc.ServeMux().Handle("/messages/v2/read_path", middleware.Build(handler.ReadPathJsonAPI))
	svc.ServeMux().Handle("/messages/v2/read_tree", middleware.Build(handler.ReadTreeJsonAPI))
	svc.ServeMux().Handle("/messages/v2/publish", middleware.Build(handler.PublishMessagesJsonAPI))
	svc.ServeMux().Handle("/messages/v2/private", middleware.Build(handler.PrivateMessagesJsonAPI))
	svc.ServeMux().Handle("/messages/v2/delete", middleware.Build(handler.DeleteMessagesJsonAPI))
	svc.ServeMux().Handle("/messages/v2/update", middleware.Build(handler.UpdateMessageJsonAPI))

	svc.ServeMux().Handle("/translations/v2/send", middleware.Build(handler.SendTranslationJsonAPI))
	svc.ServeMux().Handle("/translations/v2/update", middleware.Build(handler.UpdateTranslationJsonAPI))
	svc.ServeMux().Handle("/translations/v2/delete", middleware.Build(handler.DeleteTranslationJsonAPI))
	svc.ServeMux().Handle("/translations/v2/read", middleware.Build(handler.ReadTranslationJsonAPI))
	svc.ServeMux().Handle("/translations/v2/list", middleware.Build(handler.ListTranslationsJsonAPI))

	svc.ServeMux().Handle("POST /comments/v2/send", middleware.Build(handler.SendCommentJsonAPI))
	svc.ServeMux().Handle("POST /comments/v2/update", middleware.Build(handler.UpdateCommentJsonAPI))
	svc.ServeMux().Handle("POST /comments/v2/read", middleware.Build(handler.ReadCommentJsonAPI))
	svc.ServeMux().Handle("POST /comments/v2/delete", middleware.Build(handler.DeleteCommentsJsonAPI))
	svc.ServeMux().Handle("POST /comments/v2/list", middleware.Build(handler.ListCommentsJsonAPI))

	return nil
}
