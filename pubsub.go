package mjd

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

var ch *amqp.Channel

type PubSubData struct {
	Name        string
	Exchange    string
	RoutingKey  string
	Mandatory   bool
	Immediate   bool
	ContentType string
	Body        interface{}
}

type PubSubRabbitMq struct {
}

type PubSubInterface interface {
	Send(data PubSubData) error
	Receive(exchange string) (<-chan amqp.Delivery, error)
}

func (r *PubSubRabbitMq) Receive(exchange string) (<-chan amqp.Delivery, error) {

	err := ch.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (r *PubSubRabbitMq) Send(data PubSubData) error {
	err := ch.ExchangeDeclare(
		data.Name, // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(data.Body)
	if err != nil {
		return err
	}

	if data.ContentType == "" {
		data.ContentType = "application/json"
	}

	err = ch.PublishWithContext(ctx,
		data.Exchange,   // exchange
		data.RoutingKey, // routing key
		data.Mandatory,  // mandatory
		data.Immediate,  // immediate
		amqp.Publishing{
			ContentType: data.ContentType,
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}

func NewRabbitMq() (*PubSubRabbitMq, error) {
	config := GetConfig().Message
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", config.User, config.Pass, config.Host))
	if err != nil {
		return nil, err
	}

	ch, err = conn.Channel()
	if err != nil {
		return nil, err
	}
	return &PubSubRabbitMq{}, nil
}
