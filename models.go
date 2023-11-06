package main

import (
	"time"

	"gorm.io/gorm"
)

type student struct {
	gorm.Model
	Name       string
	Age        int
	Roll       int
	TestScores []testScore
}

type testScore struct {
	gorm.Model
	Date      time.Time
	StudentID uint
	TestName  string
	Score     string
}
