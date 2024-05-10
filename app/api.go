package app

import (
  "fmt"
  "github.com/nats-io/nats.go/micro"
)

func AttachEndpoints(srv micro.Service, options *Options, ts ThreadStore) error {
  grp := srv.AddGroup("ai")

  err := register(grp, "call", "call a model", callRequestHandler(ts, options.OllamaUrl, options.DefaultModel))
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
