package application

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type instrumentedApp struct {
	App
	messagesSaved prometheus.Counter
}

var _ App = (*instrumentedApp)(nil)

func NewInstrumentedApp(app App, messagesSaved prometheus.Counter) App {
	return instrumentedApp{
		App: app,
		messagesSaved: messagesSaved,
	}
}

func (a instrumentedApp) SaveMessage(ctx context.Context, id int64, text, title string, fileIDs []int64, userID int64, private bool, name string) (err error) {
	err = a.App.SaveMessage(ctx, id, text, title, fileIDs, userID, private, name)
	if err != nil {
		return
	}
	a.messagesSaved.Inc()
	return
}
