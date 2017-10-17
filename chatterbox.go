package chatterbox

import ()

type ChatterboxMessage struct {
	Destination string      `json:"destination"`
	Source      string      `json:"source"`
	Status      int         `json:"status"`
	Body        interface{} `json:"body"`
}

type CloudWatchMessage struct {
	Source string      `json:"source"`
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

type Receiver interface {
	Close() error
	Listen(...Dispatcher) error
}

type Dispatcher interface {
	Dispatch(ChatterboxMessage) error
}
