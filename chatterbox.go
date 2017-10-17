package chatterbox

import ()

type ChatterboxMessage struct {
	Destination string      `json:"destination"`
	Host        string      `json:"host"`
	Application string      `json:"application"`
	Context     string      `json:"context"`
	Status      string      `json:"status"`
	StatusCode  int         `json:"status_code"`
	Details     interface{} `json:"body"`
	Signature   string      `json:"signature,omitempty"`
}

type CloudWatchMessage struct {
	Host        string      `json:"host,omitempty"`
	Application string      `json:"application,omitempty"`
	Context     string      `json:"context,omitempty"`
	Status      string      `json:"status,omitempty"`
	StatusCode  int         `json:"status_code,omitempty"`
	Details     interface{} `json:"details,omitempty"`
}

type Receiver interface {
	Close() error
	Listen(...Dispatcher) error
}

type Dispatcher interface {
	Dispatch(ChatterboxMessage) error
	Close() error
}

type Broadcaster interface {
	Broadcast(ChatterboxMessage) error
	Close() error
}
