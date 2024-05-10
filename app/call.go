package app

import (
  "context"
  "github.com/henomis/lingoose/llm/ollama"
  "github.com/henomis/lingoose/thread"
  "github.com/nats-io/nats.go/micro"
  "github.com/rs/xid"
  "github.com/rs/zerolog/log"
)

const (
  ThreadIdHeader = "nats-thread-id"
  ModelHeader    = "nats-model"
)

func callRequestHandler(ts ThreadStore, endpoint string, defaultModel string) micro.Handler {
  return micro.HandlerFunc(func(req micro.Request) {
    ctx := context.Background()

    var th *thread.Thread
    tid := req.Headers().Get(ThreadIdHeader)
    if tid == "" {
      tid = xid.New().String()
      th = thread.New()
    } else {
      thrd, err := ts.GetThread(tid)
      if err != nil {
        _ = req.Error("THREAD_ERROR", err.Error(), nil)
        return
      }

      th = thrd
    }

    model := req.Headers().Get(ModelHeader)
    if model == "" {
      model = defaultModel
    }

    th.AddMessage(thread.NewUserMessage().AddContent(
      thread.NewTextContent(string(req.Data()))))

    err := ollama.New().
      WithEndpoint(endpoint).
      WithModel(model).
      WithStream(func(s string) {
        _ = req.Respond([]byte(s), micro.WithHeaders(map[string][]string{
          ThreadIdHeader: {tid},
          ModelHeader:    {model},
        }))
      }).
      Generate(ctx, th)

    if err != nil {
      _ = req.Error("LLM_ERROR", err.Error(), nil)
      return
    }

    if err := ts.StoreThread(tid, th); err != nil {
      log.Warn().Err(err).Msg("failed to store thread")
    }

  })
}
