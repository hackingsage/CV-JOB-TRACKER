package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Application struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	Company        string    `json:"company"`
	Role           string    `json:"role"`
	Status         string    `json:"status"`
	JobURL         string    `json:"jobUrl"`
	JobDescription string    `json:"jobDescription"`
	ResumeText     string    `json:"resumeText"`
	FitScore       int       `json:"fitScore"`
	Strengths      []string  `json:"strengths"`
	Gaps           []string  `json:"gaps"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
