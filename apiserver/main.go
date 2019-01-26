package main

import (
	"io"
	"log"

	"github.com/nats-io/go-nats-streaming"
)

// Convenience function to log the error on a deferred close.
func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("close error: %s", err)
	}
}

func main() {
	// Specify the cluster (of one node) and some client id for this connection.
	conn, err := stan.Connect("test-cluster", "test-client")
	if err != nil {
		log.Print(err)
		return
	}
	defer logCloser(conn)

	// Now the patterns..
	handle := func(msg *stan.Msg) {
		// If the msg is handled successfully, you can manually
		// send an ACK to the server. More importantly, if processing
		// fails, you can choose *not* send an ACK and you will receive
		// the message again later.

		// This will only fail if the connection with the server
		// has gone awry.
		if err := msg.Ack(); err != nil {
			log.Printf("failed to ACK msg: %d", msg.Sequence)
		}
	}

	sub, err := conn.Subscribe(
		"stream-name",
		handle,
		stan.SetManualAckMode(),
		stan.DurableName("i-will-remember"),
	)
	if err != nil {
		log.Print(err)
		return
	}
	defer logCloser(sub)

	<- make(chan int)
}
