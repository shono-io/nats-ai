package api

import (
	"fmt"
	"github.com/nats-io/nats.go/micro"
	"github.com/shono-io/mini"
	"nassist/app"
)

func Attach(srv *mini.Service, endpoint string, ts app.ThreadStore) error {
	grp := srv.AddGroup("ai")

	err := register(grp, "call", "call a model", CallRequestHandler(ts, endpoint))
	if err != nil {
		return err
	}

	return nil
}

func register(gr micro.Group, name, description string, handler micro.Handler) error {
	err := gr.AddEndpoint(name,
		handler,
		micro.WithEndpointMetadata(map[string]string{
			"description": description,
			"format":      "application/json",
		}))
	if err != nil {
		return fmt.Errorf("failed to add %s endpoint: %w", name, err)
	}
	return nil
}
