package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in MongoDB
type User struct {
	ID            primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"` // MongoDB ObjectID
	Username      string              `json:"username" bson:"username"`
	NmUser        string              `json:"nm_user" bson:"nm_user"`
	Pass          string              `json:"pass" bson:"pass"`
	Email         string              `json:"email" bson:"email"`
	RoleAktif     int                 `json:"role_aktif" bson:"role_aktif"`       // Hanya "admin" atau "civitas"
	IdJenisUser   int                 `json:"id_jenis_user" bson:"id_jenis_user"` // Contoh: Mahasiswa, Dosen
	CreatedAt     primitive.Timestamp `json:"created_at" bson:"created_at"`
	CreatedBy     int                 `json:"created_by" bson:"created_by"`
	UpdatedAt     primitive.Timestamp `json:"updated_at" bson:"updated_at"`
	UpdatedBy     int                 `json:"updated_by" bson:"updated_by"`
	AuthKey       string              `json:"auth_key" bson:"auth_key"`
	Photo         string              `json:"photo" bson:"photo"`
	Phone         string              `json:"phone" bson:"phone"`
	Token         string              `json:"token" bson:"token"`
	jenis_kelamin int                 `json:"jenis_kelamin" bson:"jenis_kelamin"`
}
