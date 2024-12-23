package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JenisUser struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	IDJenisUser int                `bson:"id_jenis_user" json:"id_jenis_user"`
	NmJenisUser string             `bson:"nm_jenis_user" json:"nm_jenis_user"`
	Modul       []Modul            `bson:"modul" json:"modul"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	UpdatedBy   string             `bson:"updated_by" json:"updated_by"`
}
