package mq

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type Producer struct {
	Channel     *amqp091.Channel
	Exchange    string
	RoutingKey  string
	Mandatory   bool
	Immediate   bool
	ContentType string
}

func (p *Producer) PublishWithContext(ctx context.Context, body []byte) error {
	err := p.Channel.PublishWithContext(
		ctx,
		p.Exchange,   // exchange
		p.RoutingKey, // routing key
		p.Mandatory,  // mandatory
		p.Immediate,  // immediate
		amqp091.Publishing{
			ContentType: p.ContentType,
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}
	return err
}
