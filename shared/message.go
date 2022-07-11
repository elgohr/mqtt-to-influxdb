package shared

import (
	"fmt"
	"time"
)

type Message struct {
	Topic string
	Value any
	Time  time.Time
}

func (m Message) Hash() string {
	return fmt.Sprintf("%v", m)
}
