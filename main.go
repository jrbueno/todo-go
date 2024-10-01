package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos = []Todo{
	{ID: 1, Title: "Learn Fiber", Completed: false},
	{ID: 2, Title: "Learn Golang", Completed: false},
	{ID: 3, Title: "Build RESTful API", Completed: false},
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")

	app := fiber.New()

	app.Get("api/todos", GetTodosHandler)
	app.Post("api/todos", CreateTodoHandler)
	app.Patch("api/todos/:id", updateTodoHandler)
	app.Delete("api/todos/:id", deleteTodoHandler)

	//Start Server
	log.Fatal(app.Listen(":" + PORT))
}

func deleteTodoHandler(c *fiber.Ctx) error {
	paramID := c.Params("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func updateTodoHandler(c *fiber.Ctx) error {
	type request struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}

	var body request

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	paramID := c.Params("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	var updatedTodo *Todo
	for i, todo := range todos {
		if todo.ID == id {
			if body.Title != nil {
				todos[i].Title = *body.Title
			}
			if body.Completed != nil {
				todos[i].Completed = *body.Completed
			}
			updatedTodo = &todos[i]
			break
		}
	}

	if updatedTodo == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Todo not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedTodo)
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
		ID:        len(todos) + 1,
		Title:     body.Title,
		Completed: false,
	}

	todos = append(todos, todo)

	return c.Status(fiber.StatusCreated).JSON(todo)
}

func GetTodosHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(todos)
}
