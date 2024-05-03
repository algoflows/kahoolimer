package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var quizzCollection *mongo.Collection

func main() {
	app := fiber.New()
	app.Use(cors.New())

	setupMongoDB()

	app.Get("/", index)
	app.Get("/api/quizzes", getQuizzes)

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// Send a message to the client
		err := c.WriteMessage(websocket.TextMessage, []byte("Hello from server!"))
		if err != nil {
			log.Println("Write error:", err)
			return
		}

		// Read messages from the client
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received message: %s", msg)
		}
	}))

	log.Fatal(app.Listen(":3000"))
}

func setupMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:BVFxHoSvdJEANamwrkgPtBRCigkLDJUR@viaduct.proxy.rlwy.net:53598"))
	if err != nil {
		panic(err)
	}
	quizzCollection = client.Database("test").Collection("quizzes")
}

func index(c *fiber.Ctx) error {
	return c.SendString("Hello world")
}

func getQuizzes(c *fiber.Ctx) error {
	cursor, err := quizzCollection.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	var quizzes []map[string]interface{}

	if err = cursor.All(context.Background(), &quizzes); err != nil {
		panic(err)
	}

	return c.JSON(quizzes)
}
