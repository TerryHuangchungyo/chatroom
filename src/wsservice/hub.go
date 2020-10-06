package wsservice

import (
	"time"
)

type Message struct {
	id      uint32
	time    time.Time
	content []byte
}

type Hub struct {
	Id         uint32
	Name       string
}

