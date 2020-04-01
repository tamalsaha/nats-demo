# nats-demo

```console
nats-streaming-server -m 8222 \
  --store file \
  --dir $HOME/.nats \
  --max_msgs 0 \
  --max_bytes 0
```

## nats.io basisc:
- https://docs.nats.io/nats-concepts/intro
- code examples: https://github.com/nats-io/go-nats

We are not going to use nats-streaming at this point. The new coolness is nats jetstream which is built into nats v2.0 and supports persistence. You should read this blog series to learn about using Message Oriented Middleware (MoM):

- https://choria.io/blog/post/2020/03/23/nats_patterns_1/


## nats.io streaming:
- https://docs.nats.io/nats-streaming-concepts/intro
- https://nats.io/blog/use-cases-for-persistent-logs-with-nats-streaming/

`That is, no two connections with the same client ID will be able to run concurrently.`
- https://nats-io.github.io/docs/nats_streaming/client-connections.html


## Video Explainer (audio in Bangla): https://youtu.be/J-EV3GmWkuM
