package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Product struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Price       int    `json:"price" bson:"price"`
	Description string `json:"description" bson:"description"`
}

// PublishProduct ส่งข้อมูลสินค้าไปยัง RabbitMQ
func PublishProduct(product Product) error {
	// เชื่อมต่อ RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return err
	}
	defer conn.Close()

	// เปิด channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return err
	}
	defer ch.Close()

	// สร้าง queue ชื่อ `product_queue`
	q, err := ch.QueueDeclare(
		"product_queue",
		false, // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
		return err
	}

	// แปลงข้อมูลสินค้าเป็น JSON
	body, err := json.Marshal(product)
	if err != nil {
		log.Fatalf("Failed to marshal product: %v", err)
		return err
	}

	// ส่งข้อมูลไปยัง queue
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
		return err
	}

	log.Printf(" [x] Sent product: %s", body)
	return nil
}
