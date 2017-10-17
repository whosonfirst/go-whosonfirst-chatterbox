package chatterbox

import ()

type ChatterboxMessage struct {
	Destination string      `json:"destination"`
	Host        string      `json:"host"`
	Application string      `json:"application"`
	Context     string      `json:"context"`
	Status      string      `json:"status"`
	StatusCode  int         `json:"status"`
	Details     interface{} `json:"body"`
}

type CloudWatchMessage struct {
	Host        string      `json:"host"`
	Application string      `json:"application"`
	Context     string      `json:"context"`
	Status      string      `json:"status"`
	StatusCode  int         `json:"status"`
	Details     interface{} `json:"details"`
}

type Receiver interface {
	Close() error
	Listen(...Dispatcher) error
}

type Dispatcher interface {
	Dispatch(ChatterboxMessage) error
}
