package app

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Option func(*Options)

func ViperFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("server", "s", "nats://localhost:4222", "Nats Server Url")
	cmd.PersistentFlags().String("creds", "", "Nats User Credentials")
	cmd.PersistentFlags().String("user", "", "Nats Username or Token")
	cmd.PersistentFlags().String("password", "", "Nats Password")
	cmd.PersistentFlags().String("nkey", "", "Nats User NKEY")

	cmd.PersistentFlags().String("ollama-url", "http://localhost:11434/api", "the endpoint to the ollama service")
	cmd.PersistentFlags().String("ollama-default-model", "llama3:latest", "the default model to use if none is provided")
	cmd.PersistentFlags().StringP("ttl", "t", "10m", "the amount of time to keep a thread in memory")
	cmd.PersistentFlags().IntP("count", "c", 10, "the number of threads to keep in memory")
	cmd.PersistentFlags().String("bucket", "", "kv bucket to store threads in")
	cmd.PersistentFlags().String("jsdomain", "", "kv bucket is located in js domain")

	cmd.PersistentFlags().StringP("loglevel", "l", "INFO", "the log level")
}

func FromViper(v *viper.Viper) []Option {
	var result []Option

	result = append(result, WithNatsUrl(v.GetString("server")))
	result = append(result, WithOllamaUrl(v.GetString("ollama-url")))
	result = append(result, WithLogLevel(v.GetString("loglevel")))
	result = append(result, WithDefaultModel(v.GetString("ollama-default-model")))

	// -- required
	user := v.GetString("user")
	password := v.GetString("password")
	nkey := v.GetString("nkey")
	credsFile := v.GetString("creds")

	if credsFile != "" {
		result = append(result, WithCredentialsFile(credsFile))
	} else if user != "" && nkey != "" {
		result = append(result, WithCredentials(user, nkey))
	} else if user != "" && password != "" {
		result = append(result, WithUsernamePassword(user, password))
	}

	bucket := v.GetString("bucket")
	if bucket != "" {
		result = append(result, WithKVThreadStore(bucket, v.GetString("jsdomain")))
	} else {
		result = append(result, WithMemoryThreadStore(v.GetString("ttl"), v.GetInt("count")))
	}

	return result
}

type Options struct {
	NatsUrl      string
	NatsOptions  []nats.Option
	OllamaUrl    string
	LogLevel     string
	ThreadStore  ThreadStore
	DefaultModel string
	JsDomain     string
	Bucket       string
}

func WithNatsUrl(url string) Option {
	return func(o *Options) {
		o.NatsUrl = url
	}
}

func WithUsernamePassword(username string, password string) Option {
	return func(o *Options) {
		o.NatsOptions = append(o.NatsOptions, nats.UserInfo(username, password))
	}
}

func WithCredentials(jwt string, seed string) Option {
	return func(o *Options) {
		o.NatsOptions = append(o.NatsOptions, nats.UserJWTAndSeed(jwt, seed))
	}
}

func WithCredentialsFile(file string) Option {
	return func(o *Options) {
		o.NatsOptions = append(o.NatsOptions, nats.UserCredentials(file))
	}
}

func WithNatsOptions(opts ...nats.Option) Option {
	return func(o *Options) {
		o.NatsOptions = opts
	}
}

func WithOllamaUrl(url string) Option {
	return func(o *Options) {
		o.OllamaUrl = url
	}
}

func WithMemoryThreadStore(ttl string, count int) Option {
	return func(o *Options) {
		d, _ := time.ParseDuration(ttl)
		o.ThreadStore = NewMemoryThreadStore(count, d)
	}
}

func WithKVThreadStore(bucket string, domain string) Option {
	return func(o *Options) {
		var err error
		o.ThreadStore, err = NewKvThreadStore(bucket, domain, o.NatsOptions)
		if err != nil {
			panic(err)
		}
	}
}

func WithLogLevel(level string) Option {
	return func(o *Options) {
		o.LogLevel = level
	}
}

func WithDefaultModel(model string) Option {
	return func(o *Options) {
		o.DefaultModel = model
	}
}
