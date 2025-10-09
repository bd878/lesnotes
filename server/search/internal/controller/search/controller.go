package search

type Config struct {}

type Controller struct {
	conf Config
}

func New(conf Config) *Controller {
	return &Controller{conf}
}
