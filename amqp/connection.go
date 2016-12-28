package amqp

import (
	"errors"
	"github.com/eventials/goevents/messaging"
	amqplib "github.com/streadway/amqp"
)

type Connection struct {
	connection *amqplib.Connection
}

// NewConnection returns an AMQP Connection.
func NewConnection(url string) (messaging.Connection, error) {
	conn, err := amqplib.Dial(url)

	if err != nil {
		return nil, err
	}

	return &Connection{
		conn,
	}, nil
}

func (c *Connection) NotifyConnectionClose() <-chan error {
	ch := make(chan error)

	go func() {
		ch <- errors.New((<-c.connection.NotifyClose(make(chan *amqplib.Error))).Error())
	}()

	return ch
}

// Consumer returns an AMQP Consumer.
func (c *Connection) Consumer(autoAck bool, exchange, queue string) (messaging.Consumer, error) {
	return NewConsumer(c, autoAck, exchange, queue)
}

// Producer returns an AMQP Producer.
func (c *Connection) Producer(exchange, queue string) (messaging.Producer, error) {
	return NewProducer(c, exchange, queue)
}

// Close closes the AMQP connection.
func (c *Connection) Close() {
	c.connection.Close()
}

func (c *Connection) WaitUntilConnectionCloses() {
	<-c.NotifyConnectionClose()
}