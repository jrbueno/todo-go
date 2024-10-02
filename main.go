package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title"`
	Completed bool               `json:"completed"`
}

var todoCollection *mongo.Collection

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")
	MONGOURI := os.Getenv("MONGO_URI")

	//Setop MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGOURI))
	if err != nil {
		log.Fatal(err)
	}

	//defer disconnect from MongoDB
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	//ping the MongoDB
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Error pinging MongoDB")
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")

	//Create DB & Collection
	todoCollection = client.Database("golang-db").Collection("todos")

	//Fiber Setup
	app := fiber.New()

	//Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("api/todos", GetTodosHandler)
	app.Post("api/todos", CreateTodoHandler)
	app.Patch("api/todos/:id", updateTodoHandler)
	app.Delete("api/todos/:id", deleteTodoHandler)

	//Start Server
	log.Fatal(app.Listen(":" + PORT))
}

func deleteTodoHandler(c *fiber.Ctx) error {
	paramID := c.Params("id")
	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	filter := bson.M{"_id": id}
	_, err = todoCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting todo",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func updateTodoHandler(c *fiber.Ctx) error {
	// Define request body
	type request struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}

	// Parse body
	var body request
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Get parameter id
	paramID := c.Params("id")
	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	var filter = bson.M{"_id": id}
	var todo = Todo{}
	err = todoCollection.FindOne(context.Background(), filter).Decode(&todo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding todo",
		})
	}
	// Update fields if provided in request body
	if body.Title != nil {
		todo.Title = *body.Title
	}
	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	// Save the updated todo
	_, err = todoCollection.ReplaceOne(context.Background(), filter, todo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating todo",
		})
	}

	return c.Status(fiber.StatusOK).JSON(todo)
}

func CreateTodoHandler(c *fiber.Ctx) error {
	type request struct {
		Title string `json:"title"`
	}

	var body request

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse JSON",
		})
	}

	todo := Todo{
		Title: body.Title,
	}

	result, err := todoCollection.InsertOne(context.Background(), todo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating todo",
		})
	}
	todo.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(todo)
}

func GetTodosHandler(c *fiber.Ctx) error {
	var todos []Todo
	cursor, err := todoCollection.Find(context.Background(), bson.M{}, options.Find())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding todos",
		})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo = Todo{}
		err := cursor.Decode(&todo)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error decoding todos",
			})
		}
		todos = append(todos, todo)
	}

	return c.Status(fiber.StatusOK).JSON(todos)
}
