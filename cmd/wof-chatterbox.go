package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-chatterbox"
	"gopkg.in/redis.v1"
	"log"
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

	m := chatterbox.ChatterboxMessage{
		Destination: *destination,
		Host:        *host,
		Application: *application,
		Context:     *context,
		Status:      *status,
		StatusCode:  *status_code,
	}

	msg, err := json.Marshal(m)

	if err != nil {
		log.Fatal(err)
	}

	redis_endpoint := fmt.Sprintf("%s:%d", *redis_host, *redis_port)

	redis_client := redis.NewTCPClient(&redis.Options{
		Addr: redis_endpoint,
	})

	defer redis_client.Close()

	_, err = redis_client.Ping().Result()

	if err != nil {
		log.Fatal("Failed to ping Redis server ", err)
	}

	redis_client.Publish(*redis_channel, string(msg))
}
