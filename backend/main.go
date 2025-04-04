package main

import (
	"backend/config"
	"backend/rabbitmq"
	"backend/ws"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Product struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Price       int    `json:"price" bson:"price"`
	Description string `json:"description" bson:"description"`
}

func main() {
	// เชื่อมต่อ Database
	config.ConnectMongoDB()
	config.ConnectRedis()

	// Start Consumer
	go rabbitmq.ConsumeProducts()
	go ws.HandleMessages()

	r := gin.Default()

	// ตั้งค่า CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // อนุญาตการใช้ cookies
	}))

	// API สำหรับเพิ่มสินค้า (ส่งไปยัง RabbitMQ)
	r.POST("/products", func(c *gin.Context) {
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ตรวจสอบและสร้าง id ถ้ายังไม่มี
		if product.ID == "" {
			product.ID = uuid.New().String() // auto-generate UUID
		}

		// ส่งข้อมูลไปยัง RabbitMQ
		rabbitmqProduct := rabbitmq.Product{
			ID:          product.ID,
			Name:        product.Name,
			Price:       product.Price,
			Description: product.Description,
		}
		err := rabbitmq.PublishProduct(rabbitmqProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message to RabbitMQ"})
			return
		}

		// บันทึกข้อมูลลง Redis
		productData, _ := json.Marshal(product)
		err = config.RedisClient.Set(config.Ctx, "product:"+product.ID, productData, 0).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save product to Redis"})
			return
		}

		// ไม่บันทึกข้อมูลลง MongoDB ที่นี่แล้ว
		c.JSON(http.StatusOK, gin.H{"message": "Product added successfully"})
	})

	// GET /products
	r.GET("/products", func(c *gin.Context) {
		var products []Product

		// ดึงข้อมูลจาก Redis
		keys, err := config.RedisClient.Keys(config.Ctx, "product:*").Result()
		if err == nil && len(keys) > 0 {
			for _, key := range keys {
				data, _ := config.RedisClient.Get(config.Ctx, key).Result()
				var product Product
				json.Unmarshal([]byte(data), &product)
				products = append(products, product)
			}
		} else {
			// ถ้าไม่มีข้อมูลใน Redis -> ไปดึงจาก MongoDB
			cursor, err := config.MongoDB.Collection("products").Find(config.Ctx, bson.M{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query MongoDB"})
				return
			}
			defer cursor.Close(config.Ctx)
			for cursor.Next(config.Ctx) {
				var product Product
				if err := cursor.Decode(&product); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode product"})
					return
				}
				products = append(products, product)

				// เก็บข้อมูลใหม่ใน Redis หลังจากดึงจาก MongoDB
				productData, _ := json.Marshal(product)
				err = config.RedisClient.Set(config.Ctx, "product:"+product.ID, productData, 0).Err()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save product to Redis"})
					return
				}
			}
			if err := cursor.Err(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
				return
			}
		}

		c.JSON(http.StatusOK, products)
	})

	// API สำหรับลบสินค้า
	r.DELETE("/products/:id", func(c *gin.Context) {
		productID := c.Param("id")

		// ลบข้อมูลจาก MongoDB
		_, err := config.MongoDB.Collection("products").DeleteOne(config.Ctx, bson.M{"id": productID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product from MongoDB"})
			return
		}
		log.Printf("Deleted product from MongoDB: %s", productID)

		// ลบข้อมูลจาก Redis
		err = config.RedisClient.Del(config.Ctx, "product:"+productID).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product from Redis"})
			return
		}
		log.Printf("Deleted product from Redis: %s", productID)

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	})

	// API สำหรับแก้ไขสินค้า
	r.PUT("/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// อัปเดตข้อมูลใน MongoDB
		_, err := config.MongoDB.Collection("products").UpdateOne(
			config.Ctx,
			bson.M{"id": productID},
			bson.M{"$set": product},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product in MongoDB"})
			return
		}

		// อัปเดตข้อมูลใน Redis
		productData, _ := json.Marshal(product)
		err = config.RedisClient.Set(config.Ctx, "product:"+productID, productData, 0).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product in Redis"})
			return
		}

		// ส่งข้อมูลที่อัปเดตไปยัง WebSocket
		updatedProductData, _ := json.Marshal(product)
		ws.Broadcast <- updatedProductData

		c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
	})

	// WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		ws.HandleConnections(c.Writer, c.Request)
	})

	r.Run(":8080")
}
