package controllers

import (
	"context"
	"fmt"
	"log"
	"project-crud/config"
	"project-crud/models"
	"strconv"
	"strings"
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
	id := c.Params("id") // Ambil ID dari URL
	objectID, err := isValidObjectId(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	// Parse input JSON
	modul := new(models.Modul)
	if err := c.BodyParser(modul); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	// Update data di MongoDB
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": modul}

	collection := config.GetCollection("modul")
	_, err = collection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update modul"})
	}

	return c.JSON(fiber.Map{"message": "Modul updated successfully"})
}

func DeleteModul(c *fiber.Ctx) error {
	id := c.Params("id") // Ambil ID dari URL
	objectID, err := isValidObjectId(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	// Hapus modul dari MongoDB
	filter := bson.M{"_id": objectID}

	collection := config.GetCollection("modul")
	_, err = collection.DeleteOne(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete modul"})
	}

	return c.JSON(fiber.Map{"message": "Modul deleted successfully"})
}

func AddJenisUser(c *fiber.Ctx) error {
	// Parse body request
	var jenisUser models.JenisUser
	if err := c.BodyParser(&jenisUser); err != nil {
		log.Println("Error parsing body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	// Set the created time and updated time
	jenisUser.ID = primitive.NewObjectID()
	jenisUser.IsActive = true     // default value
	jenisUser.UpdatedBy = "admin" // Set as per your system logic

	// Simpan data jenis_user ke MongoDB
	collection := config.DB.Collection("jenis_user") // Fungsi untuk mendapatkan koleksi
	_, err := collection.InsertOne(context.TODO(), jenisUser)
	if err != nil {
		log.Println("Error inserting document:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add jenis user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Jenis user added successfully",
		"data":    jenisUser,
	})
}

// Update/Delete Modul pada jenis_user
func UpdateJenisUserModul(c *fiber.Ctx) error {
	// Ambil ID dari parameter URL
	id := c.Params("id_jenis_user")
	fmt.Println("Received ID:", id)
	// Validasi format ObjectId
	objectID, err := isValidObjectId(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format", // Menampilkan error jika ID tidak valid
		})
	}

	// Parse JSON body menjadi objek JenisUser
	jenisUser := new(models.JenisUser)
	if err := c.BodyParser(jenisUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Mencari dan memperbarui data JenisUser berdasarkan ObjectId
	collection := config.GetCollection("jenis_user")
	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"nm_jenis_user": jenisUser.NmJenisUser,
			"is_active":     jenisUser.IsActive,
			"updated_by":    jenisUser.UpdatedBy,
			"modul":         jenisUser.Modul,
		},
	}

	// Melakukan update pada dokumen jenis_user
	_, err = collection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update jenis_user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "JenisUser updated successfully",
	})
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
	fmt.Println("Received User ID:", userID)

	// Validasi user ID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Ambil new_jenis_user dari query parameter
	newJenisUser := c.Query("new_jenis_user")
	fmt.Println("New Jenis User:", newJenisUser)

	// Validasi new_jenis_user
	if newJenisUser == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "new_jenis_user is required",
		})
	}

	newJenisUserInt, err := strconv.Atoi(newJenisUser)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid new_jenis_user value",
		})
	}

	// Ambil koleksi user dan jenis_user
	collectionUser := config.GetCollection("user")
	collectionJenisUser := config.GetCollection("jenis_user")

	// Reset modul untuk user
	_, err = collectionUser.UpdateOne(context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"modul": []interface{}{}}})
	if err != nil {
		fmt.Println("Error resetting user modul:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset user modul"})
	}

	// Ambil template modul dari jenis_user baru
	var jenisUserData models.JenisUser
	err = collectionJenisUser.FindOne(context.Background(), bson.M{"id_jenis_user": newJenisUserInt}).Decode(&jenisUserData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Jenis user not found",
		})
	}

	// Update user dengan modul baru dan id_jenis_user baru
	_, err = collectionUser.UpdateOne(context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"modul": jenisUserData.Modul, "id_jenis_user": newJenisUserInt}})
	if err != nil {
		fmt.Println("Error updating user modul and id_jenis_user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user modul and id_jenis_user"})
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

func isValidObjectId(id string) (primitive.ObjectID, error) {
	// Menghapus spasi tambahan jika ada
	id = strings.TrimSpace(id)

	// Mengecek panjang ID, MongoDB ObjectId harus memiliki panjang 24 karakter hex
	if len(id) != 24 {
		return primitive.NilObjectID, fmt.Errorf("Invalid ID length")
	}

	// Mengonversi string ke ObjectId
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("Invalid ObjectId format")
	}

	return objectID, nil
}
