package main

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/nats-io/go-nats-streaming"
	"github.com/tamalsaha/nats-demo/api"
	"github.com/tamalsaha/nats-demo/util"
	"io"
	"log"
)

func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("close error: %s", err)
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	<-make(chan interface{})
}

const (
	ClusterID = "test-cluster"
	ClientID = "worker-0"
)

func run() error {
	conn, err := stan.Connect(
		"test-cluster",
		ClientID,
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		return err
	}
	defer logCloser(conn)

	sub, err := conn.QueueSubscribe("create-cluster", "cluster-api-workers", func(msg *stan.Msg) {
		var info api.ClusterOperation
		err := json.Unmarshal(msg.Data, &info)
		if err != nil {
			// fail to unmarshal operation, log and ACK()
			glog.Errorf("seq = %d [redelivered = %v, data = %v, err = %v]\n", msg.Sequence, msg.Redelivered, msg.Data, err)
			return
		}

		// Look up in database to detect if this cluster id was processed before.
		// if already processed before, then do nothing, just ACK()
		// if not processed, then process now.

		util.Must(conn.Publish(info.OutputSubject, []byte("performing step 0")))
		util.DoWork()

		util.Must(conn.Publish(info.OutputSubject, []byte("performing step 1")))
		util.DoWork()

		util.Must(conn.Publish(info.OutputSubject, []byte("performing step 2")))
		util.DoWork()

		util.Must(msg.Ack())

		// Print the value and whether it was redelivered.
		glog.Infof("seq = %d [redelivered = %v]\n", msg.Sequence, msg.Redelivered)

	}, stan.SetManualAckMode() /*, stan.AckWait(time.Second)*/)
	if err != nil {
		return err
	}
	defer logCloser(sub)

	return nil
}
