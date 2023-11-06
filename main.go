package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	app := fiber.New()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&student{}, &testScore{})

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

	app.Get("/students", func(c *fiber.Ctx) error {
		var students []student
		db.Preload("TestScores").Find(&students)
		return c.JSON(students)
	})

	app.Get("/students/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var student student
		db.Preload("TestScores").First(&student, id)
		return c.JSON(student)
	})

	app.Put("/students/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var student student
		if err := db.First(&student, id).Error; err != nil {
			return c.Status(404).SendString("Student not found")
		}
		student.Name = c.FormValue("name")
		student.Age, err = strconv.Atoi(c.FormValue("age"))
		if err != nil {
			return c.Status(400).SendString("Invalid age")
		}
		student.Roll, err = strconv.Atoi(c.FormValue("roll"))
		if err != nil {
			return c.Status(400).SendString("Invalid roll")
		}
		db.Save(&student)
		return c.JSON(student)
	})

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

	port := 8080
	fmt.Printf("Listening on :%d\n", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
