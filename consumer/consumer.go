package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var telegramUrl string
var chatId string

func main() {
	LoadEnv()
	conn := CreateConn()
	defer conn.Close()
	channel := CreateChannel(conn)
	defer channel.Close()
	messages := ConsumeChannel(channel)
	ListenToChanel(messages)
}

func CreateConn() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
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

func ConsumeChannel(channel *amqp.Channel) <-chan amqp.Delivery {
	massages, err := channel.Consume(
		"FirstQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return massages
}

func ListenToChanel(massages <-chan amqp.Delivery) {

	forever := make(chan bool)
	go func() {
		fmt.Println("Received messages: ")
		for delivery := range massages {
			fmt.Printf("%s\n", string(delivery.Body))
			SendToTg(string(delivery.Body))
		}

	}()
	fmt.Println("Successfully connected to RabbitMQ instance")
	fmt.Println(" [*] - waiting for messages")
	<-forever
}

func LoadEnv() {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	telegramUrl = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
	chatId = os.Getenv("TELEGRAM_CHAT_ID")
}

func SendToTg(text string) {
	time.Sleep(5 * time.Second)
	tgResponse, err := http.PostForm(telegramUrl,
		url.Values{
			"chat_id": {chatId},
			"text":    {text},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
	}
	defer tgResponse.Body.Close()
}
