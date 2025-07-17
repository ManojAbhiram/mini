package models

import "time"

type Application struct {
	Name       string    `json:"name" bson:"name"`
	Role       string    `json:"role" bson:"role"`
	ResumePath string    `json:"resumePath" bson:"resumePath"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
}
