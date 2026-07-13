package am

type SubscriberConfig struct {
	groupName string
}

func NewSubscriberConfig(options []SubscriberOption) SubscriberConfig {
	cfg := SubscriberConfig{
		groupName: "",
	}

	for _, option := range options {
		option.configureSubscriberConfig(&cfg)
	}

	return cfg
}

type SubscriberOption interface {
	configureSubscriberConfig(*SubscriberConfig)
}

func (c SubscriberConfig) GroupName() string {
	return c.groupName
}

type GroupName string

func (n GroupName) configureSubscriberConfig(cfg *SubscriberConfig) {
	cfg.groupName = string(n)
}