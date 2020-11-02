package shared

import "time"

type Message struct {
	Topic string
	Value interface{}
	Time time.Time
}
