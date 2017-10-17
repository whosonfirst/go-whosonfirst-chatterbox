package main

import (
	"context"
	"flag"
	pubsubd "github.com/whosonfirst/go-redis-tools/pubsub"
	"github.com/whosonfirst/go-whosonfirst-chatterbox/dispatcher"
	"github.com/whosonfirst/go-whosonfirst-chatterbox/receiver"
	"log"
	"os"
	"time"
)

func start_pubsubd(ctx context.Context, cancel context.CancelFunc, host string, port int) {

	log.Println("START")

}

func main() {

	var redis_host = flag.String("redis-host", "localhost", "The Redis host to connect to.")
	var redis_port = flag.Int("redis-port", 6379, "The Redis port to connect to.")
	var redis_channel = flag.String("redis-channel", "chatterbox", "The Redis channel to publish to.")

	flag.Parse()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := pubsubd.NewServer(*redis_host, *redis_port)

	if err != nil {
		log.Fatal(err)
	}

	go func() {

		err := server.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}
	}()

	// wait for the redis server - please replace with
	// something actually better than this...

	time.Sleep(2 * time.Second)

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

	err = r.Listen(d)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
