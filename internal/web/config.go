package web

import "fmt"

type WebConfig struct {
	Host string `default:"0.0.0.0"`
	Port int    `default:"8080"`
}

func (c WebConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
