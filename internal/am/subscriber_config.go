package am

import "time"

type AckType int

const (
	AckTypeAuto AckType = iota
	AckTypeManual
)

var defaultAckWait = 5 * time.Second
var defaultMaxRedeliver = 5

type SubscriberConfig struct {
	msgFilter    []string
	groupName    string
	ackType      AckType
	ackWait      time.Duration
	maxRedeliver int
}

func NewSubscriberConfig(options []SubscriberOption) *SubscriberConfig {
	cfg := SubscriberConfig{
		msgFilter:    []string{},
		groupName:    "",
		ackType:      AckTypeAuto,
		ackWait:      defaultAckWait,
		maxRedeliver: defaultMaxRedeliver,
	}
	for _, option := range options {
		option.configureSubscriberConfig(&cfg)
	}
	return &cfg
}

type SubscriberOption interface {
	configureSubscriberConfig(*SubscriberConfig)
}

func (c *SubscriberConfig) MessageFilters() []string {
	return c.msgFilter
}

func (c *SubscriberConfig) GroupName() string {
	return c.groupName
}

func (c *SubscriberConfig) AckType() AckType {
	return c.ackType
}

func (c *SubscriberConfig) AckWait() time.Duration {
	return c.ackWait
}

func (c *SubscriberConfig) MaxRedeliver() int {
	return c.maxRedeliver
}

type MessageFilter []string

func (f MessageFilter) configureSubscriberConfig(cfg *SubscriberConfig) {
	cfg.msgFilter = f
}

type GroupName string

func (g GroupName) configureSubscriberConfig(cfg *SubscriberConfig) {
	cfg.groupName = string(g)
}

func (t AckType) configureSubscriberConfig(cfg *SubscriberConfig) {
	cfg.ackType = t
}

type AckWait time.Duration

func (a AckWait) configureSubscriberConfig(cfg *SubscriberConfig) {
	cfg.ackWait = time.Duration(a)
}

type MaxRedeliver int

func (m MaxRedeliver) configureSubscriberConfig(cfg *SubscriberConfig) {
	cfg.maxRedeliver = int(m)
}
