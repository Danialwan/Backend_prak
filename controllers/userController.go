package controllers

import (
	"context"
	"net/http"
	"project-crud/config"
	"project-crud/models"
	"project-crud/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/crypto/bcrypt"
)

var userCollection = config.GetCollection("users")

// var userCollection *mongo.Collection = config.UnairDB.Collection("users")

// Create User
// Create User

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parsing user
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&user)
	if err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	// Load timezone Asia/Jakarta untuk created_at
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	user.CreatedAt = primitive.Timestamp{T: uint32(time.Now().In(loc).Unix()), I: 0}

	// Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user.Pass = string(hashedPassword)

	// Buat user baru tanpa ID manual, MongoDB akan secara otomatis membuatkan `_id`
	newUser := bson.M{
		"username":   user.Username,
		"nm_user":    user.NmUser,
		"pass":       user.Pass,
		"email":      user.Email,
		"role_aktif": user.RoleAktif,
		"created_at": user.CreatedAt,
		// "jenis_kelamin": user.JenisKelamin,
		"photo": user.Photo,
		"phone": user.Phone,
	}

	// Insert user baru ke database
	insertResult, errIns := userCollection.InsertOne(ctx, newUser)
	if errIns != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": errIns.Error()})
	}

	// Ambil `_id` yang baru di-generate oleh MongoDB dan kembalikan dalam response
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"id":      insertResult.InsertedID, // mengambil _id dari hasil insert
	})
}

// Get All Users
func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err = cursor.All(ctx, &users); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(users)
}

func GetUsers1(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err = cursor.All(ctx, &users); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(users)
}

func GetUserByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(http.StatusOK).JSON(user)
}

func UpdateUserByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	var userUpdate models.User
	if err := c.BodyParser(&userUpdate); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updateData := bson.M{
		"username":   userUpdate.Username,
		"nm_user":    userUpdate.NmUser,
		"pass":       userUpdate.Pass,
		"email":      userUpdate.Email,
		"role_aktif": userUpdate.RoleAktif,
		// "jenis_kelamin": userUpdate.JenisKelamin,
		"photo": userUpdate.Photo,
		"phone": userUpdate.Phone,
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User updated successfully"})
}

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil input username dan password
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Cari user berdasarkan username dan password
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": loginData.Username}).Decode(&user)

	// Jika tidak ditemukan
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(loginData.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Username, user.RoleAktif, user.ID, user.RoleAktif)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Return token

	c.Locals("role", user.RoleAktif)
	c.Locals("jenis_user", user.IdJenisUser)
	// return c.JSON(fiber.Map{})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token":   token,
		"message": "Login successful!",
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil parameter ID dari URL
	userID := c.Params("id")

	// Konversi userID menjadi ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Parse request body
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Update password di database
	update := bson.M{
		"$set": bson.M{
			"pass": string(hashedPassword),
		},
	}

	// Eksekusi query update
	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Periksa apakah ada dokumen yang diupdate
	if result.ModifiedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found or password not updated"})
	}

	// Berhasil
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password updated successfully"})
}

func GetUserData(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header missing"})
	}

	// Pastikan format token valid
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	// Decode token
	token := tokenParts[1]
	payload, err := utils.DecodeJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Ambil ID dari payload (dianggap sudah ada dalam bentuk string)
	idStr, ok := payload["id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid id in token"})
	}

	// Convert ID ke ObjectID
	userID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid ObjectID format"})
	}

	// Ambil data user dari database
	user, err := getUserByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user data"})
	}

	// Return user data
	return c.JSON(fiber.Map{
		"id":       user.ID.Hex(),
		"username": user.Username,
		"role":     user.RoleAktif,
	})
}

// Fungsi untuk mengambil user berdasarkan ID
func getUserByID(userID primitive.ObjectID) (*models.User, error) {

	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
