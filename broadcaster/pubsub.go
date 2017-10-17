package broadcaster

import (
	"encoding/json"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-chatterbox"
	"gopkg.in/redis.v1"
)

type PubSubBroadcaster struct {
	chatterbox.Broadcaster
	client  *redis.Client
	options PubSubBroadcasterOptions
}

type PubSubBroadcasterOptions struct {
	Host    string
	Port    int
	Channel string
}

func NewDefaultPubSubBroadcasterOptions() PubSubBroadcasterOptions {

	opts := PubSubBroadcasterOptions{
		Host:    "localhost",
		Port:    6379,
		Channel: "chatterbox",
	}

	return opts
}

func NewPubSubBroadcaster(opts PubSubBroadcasterOptions) (chatterbox.Broadcaster, error) {

	redis_endpoint := fmt.Sprintf("%s:%d", opts.Host, opts.Port)

	redis_client := redis.NewTCPClient(&redis.Options{
		Addr: redis_endpoint,
	})

	_, err := redis_client.Ping().Result()

	if err != nil {
		return nil, err
	}

	b := PubSubBroadcaster{
		client:  redis_client,
		options: opts,
	}

	return &b, nil
}

func (b *PubSubBroadcaster) Close() error {
	return b.client.Close()
}

func (b *PubSubBroadcaster) Broadcast(m chatterbox.ChatterboxMessage) error {

	msg, err := json.Marshal(m)

	if err != nil {
		return err
	}

	// please check response, yeah?

	b.client.Publish(b.options.Channel, string(msg))
	return nil
}
