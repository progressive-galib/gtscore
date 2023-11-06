package main

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	db := InitDatabase()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello")

	})

	app.Post("/students", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		ageStr := c.FormValue("age")
		rollStr := c.FormValue("roll")
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			return c.Status(400).SendString("Invalid age")
		}
		roll, err := strconv.Atoi(rollStr)
		if err != nil {
			return c.Status(400).SendString("Invalid roll")
		}
		student := student{
			Name: name,
			Age:  age,
			Roll: roll,
		}
		db.Create(&student)
		return c.JSON(student)
	})

	// Get all students
	app.Get("/students", func(c *fiber.Ctx) error {
		var students []student
		db.Preload("TestScores").Find(&students)
		return c.JSON(students)
	})

	// Get a specific student by ID
	app.Get("/students/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var student student
		db.Preload("TestScores").First(&student, id)
		return c.JSON(student)
	})

	// Update a student's details
	app.Put("/students/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var student student
		if err := db.First(&student, id).Error; err != nil {
			return c.Status(404).SendString("Student not found")
		}
		student.Name = c.FormValue("name")
		// student.Age, err = strconv.Atoi(c.FormValue("age"))
		// if err != nil {
		// 	return c.Status(400).SendString("Invalid age")
		// }
		// student.Roll, err = strconv.Atoi(c.FormValue("roll"))
		// if err != nil {
		// 	return c.Status(400).SendString("Invalid roll")
		// }
		db.Save(&student)
		return c.JSON(student)
	})

	// Delete a student
	app.Delete("/students/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var student student
		if err := db.First(&student, id).Error; err != nil {
			return c.Status(404).SendString("Student not found")
		}
		db.Delete(&student)
		return c.SendString("Student deleted successfully")
	})

	app.Post("/students/:id/test-scores", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var student student
		if err := db.First(&student, id).Error; err != nil {
			return c.Status(404).SendString("Student not found")
		}
		testName := c.FormValue("testName")
		score := c.FormValue("score")
		dateStr := c.FormValue("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(400).SendString("Invalid date format")
		}
		test := testScore{
			StudentID: student.ID,
			TestName:  testName,
			Score:     score,
			Date:      date,
		}
		db.Create(&test)
		return c.JSON(test)
	})

	app.Get("/students/:id/test-scores", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var tests []testScore
		db.Find(&tests).Where("StudentID=?", id)
		return c.JSON(tests)

	})
}
