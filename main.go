package main

import (
	"project-crud/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routes.RouteApp(app)
	app.Listen(":3000")
}

// func main() {
// 	// Connect to the database
// 	config.ConnectDB()

// 	// Access the 'user' collection
// 	userCollection := config.GetCollection("user")

// 	// Perform a query (e.g., find all documents)
// 	cursor, err := userCollection.Find(context.TODO(), bson.M{})
// 	if err != nil {
// 		log.Fatal("Error querying user collection:", err)
// 	}
// 	defer cursor.Close(context.TODO())

// 	// Iterate through the cursor and print documents
// 	var users []bson.M
// 	if err := cursor.All(context.TODO(), &users); err != nil {
// 		log.Fatal("Error decoding user documents:", err)
// 	}

// 	for _, user := range users {
// 		fmt.Println(user)
// 	}
// }
