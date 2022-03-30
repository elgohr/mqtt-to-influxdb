package shared

import (
	"fmt"
	"time"
)

type Message struct {
	Topic string
	Value interface{}
	Time  time.Time
}

func (m Message) Hash() string {
	return fmt.Sprintf("%v", m)
}
