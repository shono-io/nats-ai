package app

import (
  "fmt"
  "github.com/nats-io/nats.go"
  "github.com/nats-io/nats.go/micro"
  "github.com/rs/zerolog"
  "github.com/rs/zerolog/log"
  "os"
  "os/signal"
  "sync"
  "syscall"
)

func Launch(opts ...Option) error {
  options := &Options{}
  for _, opt := range opts {
    opt(options)
  }

  return LaunchWithOptions(options)
}

func LaunchWithOptions(options *Options) error {
  // -- create the logger for the service
  lvl, err := zerolog.ParseLevel(options.LogLevel)
  if err != nil {
    log.Warn().Msgf("failed to parse log level, reverting to INFO")
    lvl = zerolog.InfoLevel
  }
  zerolog.SetGlobalLevel(lvl)

  nc, err := nats.Connect(options.NatsUrl, options.NatsOptions...)
  if err != nil {
    return fmt.Errorf("failed to connect to nats: %w", err)
  }

  wg := sync.WaitGroup{}
  wg.Add(1)

  scfg := micro.Config{
    Name:        "nats-ai",
    Description: "A nats service listening for AI requests",
    Version:     "0.1.0",
    DoneHandler: func(srv micro.Service) {
      wg.Done()
      log.Info().Msg("service stopped")
    },
    ErrorHandler: func(srv micro.Service, err *micro.NATSError) {
      log.Info().Str("subject", err.Subject).Err(err).Msg("service error")
    },
  }

  srv, err := micro.AddService(nc, scfg)
  if err != nil {
    return err
  }

  if err := AttachEndpoints(srv, options, options.ThreadStore); err != nil {
    return err
  }

  sig := make(chan os.Signal)

  go func() {
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
    if err := srv.Stop(); err != nil {
      log.Error().Err(err).Msg("failed to stop the service")
    }
  }()

  wg.Wait()

  return nil
}
