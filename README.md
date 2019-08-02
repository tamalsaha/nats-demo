# nats-demo

```console
nats-streaming-server -m 8222 \
  --store file \
  --dir $HOME/.nats \
  --max_msgs 0 \
  --max_bytes 0
```

## nats.io basisc:
- https://nats.io/documentation/
- https://nats.io/documentation/writing_applications/concepts/
- https://nats.io/documentation/writing_applications/subjects/
- https://nats.io/documentation/writing_applications/connecting/
- code examples: https://github.com/nats-io/go-nats

## nats.io streaming:
- https://nats.io/documentation/streaming/nats-streaming-intro/
- https://nats.io/blog/use-cases-for-persistent-logs-with-nats-streaming/

`That is, no two connections with the same client ID will be able to run concurrently.`
- https://nats-io.github.io/docs/nats_streaming/client-connections.html


## Video Explainer (audio in Bangla): https://youtu.be/J-EV3GmWkuM
