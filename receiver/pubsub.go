package receiver

import (
	"encoding/json"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-chatterbox"
	"gopkg.in/redis.v1"
	"log"
)

type PubSubReceiver struct {
	chatterbox.Receiver
	client *redis.Client
	opts   PubSubReceiverOptions
}

type PubSubReceiverOptions struct {
	Host    string
	Port    int
	Channel string
}

func NewDefaultPubSubReceiverOptions() PubSubReceiverOptions {

	opts := PubSubReceiverOptions{
		Host:    "localhost",
		Port:    6379,
		Channel: "chatterbox",
	}

	return opts
}

func NewPubSubReceiver(opts PubSubReceiverOptions) (chatterbox.Receiver, error) {

	redis_endpoint := fmt.Sprintf("%s:%d", opts.Host, opts.Port)

	redis_client := redis.NewTCPClient(&redis.Options{
		Addr: redis_endpoint,
	})

	_, err := redis_client.Ping().Result()

	if err != nil {
		return nil, err
	}

	r := PubSubReceiver{
		client: redis_client,
		opts:   opts,
	}

	return &r, nil
}

func (r *PubSubReceiver) Close() error {
	return r.client.Close()
}

func (r *PubSubReceiver) Listen(dispatchers ...chatterbox.Dispatcher) error {

	pubsub_client := r.client.PubSub()
	err := pubsub_client.Subscribe(r.opts.Channel)

	if err != nil {
		return err
	}

	for {

		i, _ := pubsub_client.Receive()

		if msg, _ := i.(*redis.Message); msg != nil {

			payload := msg.Payload
			log.Println("PAYLOAD", payload)

			var m chatterbox.ChatterboxMessage
			err := json.Unmarshal([]byte(payload), &m)

			if err != nil {
				log.Println(err)
				continue
			}

			for _, d := range dispatchers {

				err = d.Dispatch(m)

				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	return nil
}
