package app

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/henomis/lingoose/thread"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type kvThreadStore struct {
	nc  *nats.Conn
	js  jetstream.JetStream
	kv  jetstream.KeyValue
	ctx context.Context
}

func NewKvThreadStore(bucket string, domain string, opts []nats.Option) (ThreadStore, error) {
	var err error
	ts := &kvThreadStore{}

	//This is a bit awkward since there isnt any context, or nats-connection object, in the
	//bootstrapping of the app.
	ts.ctx = context.Background()

	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	ts.nc, err = nats.Connect(url, opts...)
	if err != nil {
		return nil, err
	}

	if domain != "" {
		ts.js, err = jetstream.NewWithDomain(ts.nc, domain)
	} else {
		ts.js, err = jetstream.New(ts.nc)
	}

	if err != nil {
		return nil, err
	}

	ts.kv, err = ts.js.KeyValue(ts.ctx, bucket)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (ts *kvThreadStore) GetThread(threadID string) (*thread.Thread, error) {
	ctx, cancel := context.WithTimeout(ts.ctx, time.Second*5)
	defer cancel()

	v, err := ts.kv.Get(ctx, threadID)
	if err != nil {
		return nil, err
	}
	var t thread.Thread
	if err = json.Unmarshal(v.Value(), &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (ts *kvThreadStore) StoreThread(threadID string, thread *thread.Thread) error {
	ctx, cancel := context.WithTimeout(ts.ctx, time.Second*5)
	defer cancel()

	j, err := json.Marshal(thread)
	if err != nil {
		return err
	}
	if _, err := ts.kv.Put(ctx, threadID, j); err != nil {
		return err
	}
	return nil
}
