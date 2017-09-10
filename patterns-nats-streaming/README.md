# Use cases for persistent logs with NATS Streaming

These are fully working examples of the patterns described in the post [Use cases for persistent logs with NATS Streaming]().

To run, [Go must be installed](https://golang.org/dl/). [Download NATS streaming](https://github.com/nats-io/nats-streaming-server/releases) and run it:

```
$ ./nats-streaming-server \
  --cluster_id test-cluster \
  --max_msgs 0 \
  --max_bytes 0 \
  --max_age 0s \
```

Each section belows links to the script and show example output (so you don't have to run them). If you choose to, just use `go run`, e.g. `go run 1-ephemeral.go`.

## [ephemeral](./1-ephemeral.go)

This example shows the standard subscriber.

*Note your sequence numbers may be different for all examples. They are the sequence number of the message within the log, so if you see the same sequence number twice, they are the same message.*

```
seq = 11 [redelivered = false]
seq = 12 [redelivered = false]
seq = 13 [redelivered = false]
seq = 14 [redelivered = false]
seq = 15 [redelivered = false]
seq = 16 [redelivered = false]
seq = 17 [redelivered = false]
seq = 18 [redelivered = false]
seq = 19 [redelivered = false]
seq = 20 [redelivered = false]
```

Running it again you will see the same thing but for the next 10 messages published.

## [manual ack](./2-manual-ack.go)

This example simulates failing to ack and getting a redelivery from the server. Message 34 is not acked, the remaining queued messages were processed, then 34 was redelivered a couple times until it was successful.

```
seq = 31 [redelivered = false, acked = true]
seq = 32 [redelivered = false, acked = true]
seq = 33 [redelivered = false, acked = true]
seq = 34 [redelivered = false, acked = false]
seq = 35 [redelivered = false, acked = true]
seq = 36 [redelivered = false, acked = true]
seq = 37 [redelivered = false, acked = true]
seq = 38 [redelivered = false, acked = true]
seq = 39 [redelivered = false, acked = true]
seq = 40 [redelivered = false, acked = true]
seq = 34 [redelivered = true, acked = false]
seq = 34 [redelivered = true, acked = false]
seq = 34 [redelivered = true, acked = true]
```

## [durable](./3-durable.go)

With this example, a subscriber is *disconnected* after 5 messages, reconnects and processes the remaining messages.

```
seq = 21 [redelivered = false]
seq = 22 [redelivered = false]
seq = 23 [redelivered = false]
seq = 24 [redelivered = false]
seq = 25 [redelivered = false]
2017/09/10 16:04:04 subscriber disconnected..
2017/09/10 16:04:04 reconnecting..
seq = 26 [redelivered = true]
seq = 27 [redelivered = false]
seq = 28 [redelivered = false]
seq = 29 [redelivered = false]
seq = 30 [redelivered = false]
```

## [exactly once](./4-exactly-once.go)

This example emulates a failed ACK, but when the a subscription receives the same msg, it does not reprocess it. It only acknowledges it.

```
seq = 81 [redelivered = false, acked = true, processed = true]
seq = 82 [redelivered = false, acked = true, processed = true]
seq = 83 [redelivered = false, acked = true, processed = true]
seq = 84 [redelivered = false, acked = true, processed = true]
seq = 85 [redelivered = false, acked = true, processed = true]
seq = 86 [redelivered = false, acked = false, processed = true]
seq = 87 [redelivered = false, acked = false, processed = true]
seq = 88 [redelivered = false, acked = false, processed = true]
seq = 89 [redelivered = false, acked = false, processed = true]
seq = 90 [redelivered = false, acked = true, processed = true]
seq = 86 [redelivered = true, acked = true, processed = false]
seq = 87 [redelivered = true, acked = true, processed = false]
seq = 88 [redelivered = true, acked = true, processed = false]
seq = 89 [redelivered = true, acked = true, processed = false]
```
