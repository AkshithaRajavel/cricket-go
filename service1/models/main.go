package models

import (
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
	Country       string             `json:"country" bson:"country,omitempty"`
	Captain       bool               `json:"captain" bson:"captain,omitempty"`
}
type Team struct {
	Country string
	Players []Player
	Captain Player
}
