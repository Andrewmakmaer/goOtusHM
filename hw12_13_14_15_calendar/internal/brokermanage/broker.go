package brokermanage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	channel *amqp.Channel
	q       amqp.Queue
}

func NewBroker(endpoint, queue string) (Broker, error) {
	conn, err := amqp.Dial(endpoint)
	if err != nil {
		return Broker{}, err
	}

	var brocker Broker
	brocker.channel, err = conn.Channel()
	if err != nil {
		return Broker{}, err
	}

	brocker.q, err = brocker.channel.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return Broker{}, err
	}
	return brocker, nil
}

func (b *Broker) SendMessage(event storage.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	notification := storage.NewNotification(event.ID, event.Title, event.StartTime, event.UserID)
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	err = b.channel.PublishWithContext(ctx, "",
		b.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}
	return nil
}

func (b *Broker) ReadMessages() (<-chan amqp.Delivery, error) {
	message, err := b.channel.Consume(
		b.q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return message, nil
}
