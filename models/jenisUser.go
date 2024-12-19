package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JenisUser struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Modul []Modul            `json:"modul"`
}
