# NATS AI
This project allows you to run a NATS service which exposes a Micro service that can be used to communicate with LLM's

## Getting Started
Start by pulling the `llama3` model into ollama:
```shell
ollama pull llama3
```

Start the application:
```shell
nats-ai
```

Call the endpoint
```shell
nats req --replies=0 'ai.call' 'Hi there, Can you generate me some Benthos code?'
```

Many replies will be sent back to you, so you need to set the number of replies to `0` and depend on
the reply-timeout (defaults to `300ms`) to know when to stop listening for replies.

Each reply will hold the following headers:
- `nats-model`: the model which was used to generate the reply
- `nats-thread-id`: the thread id to identify the conversation thread 

## Using a different model
By default the `llama3` model is being used as the model of choice, but this can be
overwritten using request headers:
```shell
nats req --replies=0 --header='model:deepseek-coder:33b' 'ai.call' 'Hi there, Can you generate me some Benthos code?'
```

## Resuming the conversation
Since each response returns a `nats-thread-id` header, you can use this to resume the conversation:
```shell
nats req --replies=0 --header='nats-thread-id:1234' 'ai.call' 'Can you explain this as if it was to a 2 year old?'
```