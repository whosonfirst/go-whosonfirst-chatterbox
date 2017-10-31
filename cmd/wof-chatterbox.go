package main

import (
	"encoding/json"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-chatterbox"
	"github.com/whosonfirst/go-whosonfirst-chatterbox/broadcaster"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var redis_host = flag.String("redis-host", "localhost", "The Redis host to connect to.")
	var redis_port = flag.Int("redis-port", 6379, "The Redis port to connect to.")
	var redis_channel = flag.String("redis-channel", "chatterbox", "The Redis channel to publish to.")

	var destination = flag.String("destination", "", "")
	var host = flag.String("host", "", "")
	var application = flag.String("application", "", "")
	var context = flag.String("context", "", "")
	var status = flag.String("status", "", "")
	var status_code = flag.Int("status-code", 0, "")

	flag.Parse()

	if *destination == "" {
		log.Fatal("missing destination")
	}

	if *host == "" {

		hostname, err := os.Hostname()

		if err != nil {
			log.Fatal(err)
		}

		*host = hostname
	}

	if *application == "" {

		abs_path, err := filepath.Abs(os.Args[0])

		if err != nil {
			log.Fatal(err)
		}

		*application = filepath.Base(abs_path)
	}

	if *context == "" {
		*context = "chatterbox"
	}

	if *status == "" {
		*status = "ok"
	}

	if *status_code == 0 {
		*status_code = 1
	}

	m := chatterbox.ChatterboxMessage{
		Destination: *destination,
		Host:        *host,
		Application: *application,
		Context:     *context,
		Status:      *status,
		StatusCode:  *status_code,
	}

	var details interface{}

	args := flag.Args()

	body := strings.Join(args, " ")
	body = strings.Trim(body, " ")

	if body != "" {

		var stub interface{}

		err := json.Unmarshal([]byte(body), &stub)

		if err == nil {
			details = stub
		} else {
			details = body
		}

		m.Details = details
	}

	opts := broadcaster.NewDefaultPubSubBroadcasterOptions()
	opts.Host = *redis_host
	opts.Port = *redis_port
	opts.Channel = *redis_channel

	br, err := broadcaster.NewPubSubBroadcaster(opts)

	if err != nil {
		log.Fatal(err)
	}

	defer br.Close()

	err = br.Broadcast(m)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
