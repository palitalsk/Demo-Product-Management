package rabbitmq

import (
	"backend/config"
	"backend/ws"
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// รับข้อมูลจาก RabbitMQ และบันทึกลง MongoDB & Redis
func ConsumeProducts() {
	// เชื่อมต่อ RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// สร้าง queue
	q, err := ch.QueueDeclare("product_customer", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// รับ message จาก queue
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// รับและประมวลผล message
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var product Product
			err := json.Unmarshal(d.Body, &product)
			if err != nil {
				log.Printf("Error decoding message: %v", err)
				continue
			}

			// บันทึกลง MongoDB
			_, err = config.MongoDB.Collection("products").InsertOne(context.TODO(), product)
			if err != nil {
				log.Printf("Error inserting to MongoDB: %v", err)
			} else {
				log.Printf("Inserted into MongoDB: %+v", product)
			}

			// บันทึกลง Redis
			jsonData, _ := json.Marshal(product)
			err = config.RedisClient.Set(context.TODO(), "product:"+product.ID, jsonData, 0).Err()
			if err != nil {
				log.Printf("Error inserting to Redis: %v", err)
			} else {
				log.Printf("Inserted into Redis: %+v", product)
			}

			// ส่งข้อมูลไปยัง WebSocket
			ws.Broadcast <- jsonData
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
