package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var quizzCollection *mongo.Collection

func main() {
	godotenv.Load()                     // load the .env file
	app := fiber.New()                  // create a new fiber app
	app.Use(cors.New())                 // cors middleware to allow all origins
	app.Use("/ws", webSocketMiddleware) // websocket middleware to check if the request is a websocket connection
	connectMongoDB()                    // creates the connection to the database

	app.Get("/", index)
	app.Get("/api/quizzes", getQuizzes)
	app.Get("/ws/:id", getWebsocketId())

	log.Fatal(app.Listen(":3000"))
}

func webSocketMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func connectMongoDB() {
	uri := os.Getenv("MONGO_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
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

func getWebsocketId() func(c *fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
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
	})
}
