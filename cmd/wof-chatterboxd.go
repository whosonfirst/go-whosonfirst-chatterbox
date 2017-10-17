package main

import (
	"flag"
	pubsubd "github.com/whosonfirst/go-redis-tools/pubsub"
	"github.com/whosonfirst/go-whosonfirst-chatterbox/dispatcher"
	"github.com/whosonfirst/go-whosonfirst-chatterbox/receiver"
	"log"
	"os"
)

func main() {

	var redis_host = flag.String("redis-host", "localhost", "The Redis host to connect to.")
	var redis_port = flag.Int("redis-port", 6379, "The Redis port to connect to.")
	var redis_channel = flag.String("redis-channel", "chatterbox", "The Redis channel to publish to.")

	flag.Parse()

	server, err := pubsubd.NewServer(*redis_host, *redis_port)

	if err != nil {
		log.Fatal(err)
	}

	ready := make(chan bool)

	go func() {

		err := server.ListenAndServeWithReadySignal(ready)

		if err != nil {
			log.Fatal(err)
		}
	}()

	sig := <-ready

	if !sig {
		log.Fatal("Received negative ready signal from PubSub server")
	}

	opts := receiver.NewDefaultPubSubReceiverOptions()
	opts.Host = *redis_host
	opts.Port = *redis_port
	opts.Channel = *redis_channel

	r, err := receiver.NewPubSubReceiver(opts)

	if err != nil {
		log.Fatal(err)
	}

	defer r.Close()

	d, err := dispatcher.NewCloudWatchDispatcher()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening for chatterbox requests on the %s channel at %s:%d\n", *redis_channel, *redis_host, *redis_port)

	err = r.Listen(d)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
