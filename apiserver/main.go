package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync/atomic"
	"time"

	"github.com/go-macaron/binding"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/tamalsaha/nats-demo/api"
	macaron "gopkg.in/macaron.v1"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// ref: https://github.com/square/go-jose/blob/v2/jwt/example_test.go
// Use an RSA private key to sign
var sharedKey = []byte("secret")

func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("close error: %s", err)
	}
}

const (
	ClusterID = "test-cluster"
	ClientID  = "apiserver-0"
)

var clusterId int64 = 0

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
	m.Get("/", func() string {
		return "Hello world!"
	})

	m.Post("/token", func(ctx *macaron.Context) {
		sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: sharedKey}, (&jose.SignerOptions{}).WithType("JWT"))
		if err != nil {
			ctx.Error(500, err.Error())
			return
		}

		now := time.Now()
		cl := jwt.Claims{
			Subject:   "subject",
			Issuer:    "issuer",
			NotBefore: jwt.NewNumericDate(now),
			Expiry:    jwt.NewNumericDate(now.Add(30 * 24 * time.Hour)),
			Audience:  jwt.Audience{"ws"},
		}
		raw, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
		if err != nil {
			ctx.Error(500, err.Error())
			return
		}
		ctx.JSON(201, map[string]string{
			"token": raw,
		})
	}) // require authentication
	m.Post("/verify_token", binding.Bind(api.TokenForm{}), func(ctx *macaron.Context, data api.TokenForm) {
		token, err := jwt.ParseSigned(data.Token)
		if err != nil {
			ctx.Error(500, err.Error())
			return
		}

		out := jwt.Claims{}
		if err := token.Claims(sharedKey, &out); err != nil {
			ctx.Error(401, err.Error()) // ask for a new token
			return
		}
		ctx.JSON(200, out)
	})

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
