package main

import (
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/cryptide-dev/noise"
	"github.com/cryptide-dev/noise/kademlia"
)

func main() {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(zap.PanicLevel))
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	node, err := noise.NewNode(noise.WithNodeLogger(logger), noise.WithNodeBindPort(9000))
	if err != nil {
		panic(err)
	}
	defer node.Close()

	overlay := kademlia.New()
	node.Bind(overlay.Protocol())

	if err := node.Listen(); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
