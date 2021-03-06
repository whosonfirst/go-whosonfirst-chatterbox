package dispatcher

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/eltorocorp/cloudwatch"
	"github.com/whosonfirst/go-whosonfirst-chatterbox"
	_ "log"
	"strings"
)

type CloudWatchDispatcher struct {
	chatterbox.Dispatcher
	client *cloudwatchlogs.CloudWatchLogs
}

func NewCloudWatchDispatcher() (chatterbox.Dispatcher, error) {

	// please make me config flags and options and
	// all that good stuff (20171017/thisisaaronland)

	cfg := aws.NewConfig()
	cfg.WithRegion("us-east-1")

	sess := session.New(cfg)
	client := cloudwatchlogs.New(sess)

	d := CloudWatchDispatcher{
		client: client,
	}

	return &d, nil
}

func (d *CloudWatchDispatcher) Close() error {
	return nil
}

func (d *CloudWatchDispatcher) Dispatch(m chatterbox.ChatterboxMessage) error {

	// please for to be error checking...

	dest := strings.Split(m.Destination, "#")

	if len(dest) != 2 {
		return errors.New("Invalid destination string")
	}

	group := dest[0]
	stream := dest[1]

	gr, err := cloudwatch.AttachGroup(group, d.client)

	if err != nil {
		return err
	}

	wr, err := gr.AttachStream(stream)

	if err != nil {
		return err
	}

	cw := chatterbox.CloudWatchMessage{
		Host:        m.Host,
		Application: m.Application,
		Context:     m.Context,
		Status:      m.Status,
		StatusCode:  m.StatusCode,
		Details:     m.Details,
	}

	enc, err := json.Marshal(cw)

	if err != nil {
		return err
	}

	_, err = wr.Write(enc)

	if err != nil {
		return err
	}

	// cloudwatch library does a whole flush-on-a-timer thing

	return nil
}
