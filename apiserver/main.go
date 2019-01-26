package main

import (
	"encoding/json"
	"fmt"
	"github.com/tamalsaha/nats-demo/api"
	"io"
	"log"
	"gopkg.in/macaron.v1"
	stan "github.com/nats-io/go-nats-streaming"
	"sync/atomic"
)

func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("close error: %s", err)
	}
}

const (
	ClusterID = "test-cluster"
	ClientID = "apiserver-0"
)

func main() {
	conn, err := stan.Connect(
		ClusterID,
		ClientID,
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer logCloser(conn)

	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Renderer())

	var clusterId int64 = 0

	m.Post("/clusters", func(ctx *macaron.Context) {
		id := atomic.AddInt64(&clusterId, 1) // should be a sequence from database
		out := fmt.Sprintf("cluster-%d", id)

		op := api.ClusterOperation{
			ClusterId:     id,
			OutputSubject: out,
		}
		data, err := json.Marshal(op)
		if err != nil {
			ctx.Error(500, err.Error())
			return
		}

		err = conn.Publish("create-cluster", data)
		if err != nil {
			ctx.Error(500, err.Error())
			return
		}

		ctx.JSON(201, api.ClusterCreateResponse{OutputChannel: out})
	})
	m.Run()
}
