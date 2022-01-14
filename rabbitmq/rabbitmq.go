/*
@Author : yaokun
@Time : 2020/8/14 11:11
*/

package main

import (
	_ "context"
	"fmt"
	"github.com/streadway/amqp"
)

//Rabbitmq 初始化rabbitmq连接
type Rabbitmq struct {
	conn *amqp.Connection
	err  error
}

func New(ip string) (*Rabbitmq, error) {
	amqps := fmt.Sprintf("amqp://testmq:testmq@%s:5672/", ip)
	conn, err := amqp.Dial(amqps)
	if err != nil {
		return nil, err
	}
	rabbitmq := &Rabbitmq{
		conn: conn,
	}
	return rabbitmq, nil
}

func (rabbitmq *Rabbitmq) CreateQueue(id string) error {
	ch, err := rabbitmq.conn.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		id,    // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	return nil
}


func (rabbitmq *Rabbitmq) ClearQueue(id string) (string, error) {
	ch, err := rabbitmq.conn.Channel()
	defer ch.Close()
	if err != nil {
		return "", err
	}
	_, err = ch.QueuePurge(id, false)
	if err != nil {
		return "", err
	}
	return "Clear queue success", nil
}

