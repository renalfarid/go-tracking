package client

import (
	"context"
	"errors"
	"fmt"

	"go-tracking/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueName = "tracking_status"
)

type rabbitmqClient struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	connString    string
	trackingStatus <-chan amqp.Delivery
}

func NewRabbitMQClient(connectionString string) (*rabbitmqClient, error) {
	c := &rabbitmqClient{}
	var err error

	c.conn, err = amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	err = c.configureQueue()

	return c, err
}

func (c *rabbitmqClient) Consume(ctx context.Context) ([]byte, error) {
	for msg := range c.trackingStatus {
			_ = msg.Ack(false)
			return msg.Body, nil
		
	}
	return nil, errors.New("err when getting tracking status on channel")
}

func (c *rabbitmqClient) Publish(p models.Tracking) {
	jsonStr := fmt.Sprintf(`{ "id": %q, "longitude": %f, "latitude": %f }`, p.Id, p.Longitude, p.Latitude)

	_ = c.ch.Publish("", QueueName, true, false, amqp.Publishing{
		ContentType: "application/json",
		MessageId:   p.Id,
		Body:        []byte(jsonStr),
	})
}

func (c *rabbitmqClient) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *rabbitmqClient) configureQueue() error {
	_, err := c.ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.trackingStatus, err = c.ch.Consume(
		QueueName,
		"",
		false,
		false,
		false,
		true,
		nil,
	)
	fmt.Println(c.trackingStatus)
	return err
}