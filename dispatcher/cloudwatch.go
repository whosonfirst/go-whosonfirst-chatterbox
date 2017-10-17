package dispatcher

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/ejholmes/cloudwatch"
	"github.com/whosonfirst/go-whosonfirst-chatterbox"
	"log"
	"strings"
)

type CloudWatchDispatcher struct {
	chatterbox.Dispatcher
	client *cloudwatchlogs.CloudWatchLogs
}

func NewCloudWatchDispatcher() (chatterbox.Dispatcher, error) {

	cfg := aws.NewConfig()
	cfg.WithRegion("us-east-1")

	sess := session.New(cfg)
	client := cloudwatchlogs.New(sess)

	d := CloudWatchDispatcher{
		client: client,
	}

	return &d, nil
}

func (d *CloudWatchDispatcher) Dispatch(m chatterbox.ChatterboxMessage) error {

     	// please for to be error checking...

	dest := strings.Split(m.Destination, "#")
	group := dest[0]
	stream := dest[1]

	wr := cloudwatch.NewWriter(group, stream, d.client)

	cw := chatterbox.CloudWatchMessage{
		Source: m.Source,
		Status: m.Status,
		Body:   m.Body,
	}

	enc, err := json.Marshal(cw)

	if err != nil {
		return err
	}

	log.Println(dest, string(enc))

	_, err = wr.Write(enc)

	if err != nil {
		return err
	}

	/*
	err = wr.Flush()

	if err != nil {
		return err
	}
	*/

	return nil
}
