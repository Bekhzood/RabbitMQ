package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

const (
	low = iota
	medium
	high
)

func main() {
	conn := CreateConn()
	defer conn.Close()
	channel := CreateChannel(conn)
	defer channel.Close()
	DeclareQueue(channel)
	Send(channel)
}

func CreateConn() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Successfully connected to RabbitMQ instance")
	return conn
}

func CreateChannel(conn *amqp.Connection) *amqp.Channel {
	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return channel
}

func DeclareQueue(channel *amqp.Channel) {
	_, err := channel.QueueDeclare("FirstQueue", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func Send(channel *amqp.Channel) {
	err := channel.Publish(
		"",
		"FirstQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/pain",
			Body:        []byte("hello world high 1"),
			Priority:    high,
		},
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Successfully published message to queue")
}
