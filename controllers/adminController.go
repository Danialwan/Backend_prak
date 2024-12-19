package controllers

import (
	"context"
	"project-crud/config"
	"project-crud/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CRUD modul
func CreateModul(c *fiber.Ctx) error {
	modul := new(models.Modul)
	if err := c.BodyParser(modul); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse body"})
	}

	modul.ID = primitive.NewObjectID()
	modul.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}
	modul.UpdatedAt = modul.CreatedAt

	collection := config.GetCollection("modul")
	_, err := collection.InsertOne(context.Background(), modul)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create modul"})
	}

	return c.Status(fiber.StatusOK).JSON(modul)
}

func UpdateModul(c *fiber.Ctx) error {
	id := c.Params("id")
	modulID, _ := primitive.ObjectIDFromHex(id)

	var updateData models.Modul
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	updateData.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}

	collection := config.GetCollection("modul")
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": modulID}, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update modul"})
	}

	return c.JSON(fiber.Map{"message": "Modul updated successfully"})
}

func DeleteModul(c *fiber.Ctx) error {
	id := c.Params("id")
	modulID, _ := primitive.ObjectIDFromHex(id)

	collection := config.GetCollection("modul")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": modulID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete modul"})
	}

	return c.JSON(fiber.Map{"message": "Modul deleted successfully"})
}

// Update/Delete Modul pada jenis_user
func UpdateJenisUserModul(c *fiber.Ctx) error {
	idJenisUser := c.Params("id_jenis_user")
	var modulData models.Modul

	if err := c.BodyParser(&modulData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	collection := config.GetCollection("jenis_user")
	_, err := collection.UpdateOne(context.Background(),
		bson.M{"id_jenis_user": idJenisUser},
		bson.M{"$addToSet": bson.M{"modul": modulData}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update jenis_user modul"})
	}

	return c.JSON(fiber.Map{"message": "Modul added to jenis_user successfully"})
}

func DeleteJenisUserModul(c *fiber.Ctx) error {
	idJenisUser := c.Params("id_jenis_user")
	modulID := c.Query("modul_id")

	collection := config.GetCollection("jenis_user")
	_, err := collection.UpdateOne(context.Background(),
		bson.M{"id_jenis_user": idJenisUser},
		bson.M{"$pull": bson.M{"modul": bson.M{"id": modulID}}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete modul from jenis_user"})
	}

	return c.JSON(fiber.Map{"message": "Modul removed from jenis_user successfully"})
}

// Fungsi Pindah Jenis User
func PindahJenisUser(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	newJenisUser := c.Query("new_jenis_user")

	collectionUser := config.GetCollection("user")
	collectionJenisUser := config.GetCollection("jenis_user")

	// Hapus modul lama user
	_, err := collectionUser.UpdateOne(context.Background(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"modul": []interface{}{}}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset user modul"})
	}

	// Ambil template_modul dari jenis_user baru
	var jenisUserData models.JenisUser
	err = collectionJenisUser.FindOne(context.Background(), bson.M{"id_jenis_user": newJenisUser}).Decode(&jenisUserData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Jenis user not found"})
	}

	// Update user dengan modul baru
	_, err = collectionUser.UpdateOne(context.Background(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"modul": jenisUserData.Modul}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user modul"})
	}

	return c.JSON(fiber.Map{"message": "User updated to new jenis_user successfully"})
}

// CUD Modul Khusus User
func AddModulToUser(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	var modulData models.Modul

	if err := c.BodyParser(&modulData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	collection := config.GetCollection("user")
	_, err := collection.UpdateOne(context.Background(),
		bson.M{"_id": userID},
		bson.M{"$addToSet": bson.M{"modul": modulData}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add modul to user"})
	}

	return c.JSON(fiber.Map{"message": "Modul added to user successfully"})
}
