package queue

import (
	"fmt"
	"log"
	"reflect"
)

const (
	RabbitMQ QueueType = iota
)

type QueueType int

func New(qt QueueType, cfg any) (q *Queue, err error) {
	rt := reflect.TypeOf(cfg)
	switch qt {
	case RabbitMQ:
		if rt.Name() != "RabbitMQConfig" {
			return nil, fmt.Errorf("config need'' to be of type RabbitMQ")
		}
		conn, err := newRabbitConnection(cfg.(RabbitMQConfig))
		if err != nil {
			return nil, err
		}

		q.qc = conn
	default:
		log.Println("type not implemented yet")
	}
	return
}

type QueueConnection interface {
	Publish(body []byte) error
	Consume(chan<- QueueDto) error
}

type Queue struct {
	qc QueueConnection
}

func (q *Queue) Publish(body []byte) error {
	return q.qc.Publish(body)
}

func (q *Queue) Consume(cdto chan<- QueueDto) error {
	return q.qc.Consume(cdto)
}
