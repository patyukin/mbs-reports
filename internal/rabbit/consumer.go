package rabbit

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConfirmationCodeMessage struct {
	UserID int64  `json:"user_id"`
	Code   string `json:"code"`
}

func StartConsumer(rabbitMQURL string, service *application.NotificationService) error {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"confirmation_code_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var msg ConfirmationCodeMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("Ошибка парсинга сообщения: %v", err)
				continue
			}

			if err := service.SendConfirmationCode(msg.UserID, msg.Code); err != nil {
				log.Printf("Ошибка отправки кода подтверждения: %v", err)
			}
		}
	}()

	log.Println("RabbitMQ консьюмер запущен")
	return nil
}
