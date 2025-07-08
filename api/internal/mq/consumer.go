package mq

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	Channel   *amqp091.Channel
	QueueName string
	WorkerNum int
	AutoAck   bool
	Handler   func([]byte) error
}

func (c *Consumer) Start() {
	msgs, err := c.Channel.Consume(
		c.QueueName,
		"",
		c.AutoAck,
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	for i := 0; i < c.WorkerNum; i++ {
		go func(workerID int) {
			for d := range msgs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("[Worker %d] Recovered in handler: %v", workerID, r)
						}
					}()
					err := c.Handler(d.Body)
					if err != nil {
						log.Printf("[Worker %d] Handler error: %v", workerID, err)
					}
					if !c.AutoAck {
						// 手动 ack
						if err == nil {
							d.Ack(false)
						} else {
							d.Nack(false, true) // requeue
						}
					}
				}()
			}
		}(i)
	}
}
