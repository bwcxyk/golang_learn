/*
@Author : YaoKun
@Time : 2023/5/9 14:58
*/

package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {
	// 连接到本地的RabbitMQ服务器
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@192.168.0.11/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 创建一个通道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 创建一个名为"hello"的队列，如果不存在的话
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 发送一条消息到"hello"队列
	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
