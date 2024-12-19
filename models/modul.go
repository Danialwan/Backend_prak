package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Modul struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NmModul    string             `json:"nm_modul"`
	KetModul   string             `json:"ket_modul"`
	KategoriID int                `json:"kategori_id"`
	IsAktif    string             `json:"is_aktif"`
	Alamat     string             `json:"alamat"`
	Urutan     int                `json:"urutan"`
	GbrIcon    string             `json:"gbr_icon"`
	CreatedAt  primitive.Timestamp
	UpdatedAt  primitive.Timestamp
}
