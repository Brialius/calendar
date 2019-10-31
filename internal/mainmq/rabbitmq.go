package mainmq

import (
	"context"
	"encoding/json"
	"github.com/Brialius/calendar/internal/domain/models"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMq struct {
	ch   *amqp.Channel
	conn *amqp.Connection
}

func NewRabbitMqQueue(url string) (*RabbitMq, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMq{ch: ch, conn: conn}, nil
}

func (r *RabbitMq) DeclareQueue(ctx context.Context, qName string, durable bool) error {
	_, err := r.ch.QueueDeclare(
		qName,
		durable,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (r *RabbitMq) BindQueue(ctx context.Context, qName, routingKey, exchange string, durable bool) error {
	err := r.ch.QueueBind(
		qName,
		routingKey,
		exchange,
		false,
		nil,
	)
	return err
}

func (r *RabbitMq) DeclareExchange(ctx context.Context, name, kind string, durable bool) error {
	err := r.ch.ExchangeDeclare(
		name,
		kind,
		durable,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (r *RabbitMq) SetQos(ctx context.Context, prefetchCount, prefetchSize int, global bool) error {
	return r.ch.Qos(
		prefetchCount,
		prefetchSize,
		global,
	)
}

func (r *RabbitMq) SendTaskToQueue(ctx context.Context, exchange, routingKey string, event *models.Event) error {
	jsonBody, err := json.Marshal(event)
	if err != nil {
		log.Printf("can't marshal to JSON  `%v`: %s", event, err)
		return err
	}
	return r.ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		})
}

func (r *RabbitMq) ConsumeTasksFromQueue(ctx context.Context, qName, consumer string, autoAck bool, task func(ctx context.Context, event *models.Event) error) error {
	msgs, err := r.ch.Consume(
		qName,
		consumer,
		autoAck,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			e := &models.Event{}
			log.Printf("Received a message: %s", d.Body)
			err = json.Unmarshal(d.Body, e)
			if err != nil {
				log.Printf("can't unmarshal from JSON  `%v`: %s", e, err)
			}
			if task(ctx, e) == nil {
				_ = d.Ack(false)
			}
		}
	}()
	<-forever

	return nil
}

func (r *RabbitMq) Close(ctx context.Context) {
	_ = r.ch.Close()
	_ = r.conn.Close()
}
