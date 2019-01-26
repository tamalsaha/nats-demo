package util

import (
	"math/rand"
	"runtime/debug"
	"time"
)

func Must(err error) {
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
}

func DoWork() {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}
