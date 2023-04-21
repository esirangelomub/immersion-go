package queue

import (
	"context"
	ampq "github.com/rabbitmq/amqp091-go"
	"time"
)

type RabbitMQConfig struct {
	URL       string
	TopicName string
	Timeout   time.Duration
}

type RabbitConnection struct {
	cfg  RabbitMQConfig
	conn *ampq.Connection
}

func newRabbitConnection(cfg RabbitMQConfig) (*RabbitConnection, error) {
	conn, err := ampq.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	return &RabbitConnection{
		cfg:  cfg,
		conn: conn,
	}, nil
}

func (rc *RabbitConnection) Publish(body []byte) error {
	c, err := rc.conn.Channel()
	if err != nil {
		return err
	}

	mp := ampq.Publishing{
		DeliveryMode: ampq.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         body,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30)
	defer cancel()
	return c.PublishWithContext(
		ctx,
		"",
		rc.cfg.TopicName,
		false,
		false,
		mp,
	)
}

func (rc *RabbitConnection) Consume(cdto chan<- QueueDto) error {
	ch, err := rc.conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(rc.cfg.TopicName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for d := range msgs {
		dto := QueueDto{}
		dto.FromJson(d.Body)

		cdto <- dto
	}

	return nil
}
