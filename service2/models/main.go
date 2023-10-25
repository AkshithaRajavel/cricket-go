package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	ID            primitive.ObjectID `json:"_id,omitempty" required:"true" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name,omitempty"`
	Jersey        string             `json:"Jersey" bson:"Jersey,omitempty"`
	Age           int                `json:"age" bson:"age,omitempty"`
	PrimaryRole   string             `json:"primary_role" bson:"primary_role,omitempty"`
	SecondaryRole []string           `json:"secondary_role" bson:"secondary_role,omitempty"`
	Matches       int                `json:"matches" bson:"matches,omitempty"`
	Runs          int                `json:"runs" bson:"runs,omitempty"`
	Wickets       int                `json:"wickets" bson:"wickets,omitempty"`
}
type Match struct {
	Name   string
	Team1  string
	Team2  string
	Winner string
	Date   time.Time
}
