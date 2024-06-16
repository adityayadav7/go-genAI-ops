package main

import (
	"fmt"
	"log"
	"net/http"
	"crypto/rand"
    "encoding/base64"

	"crudapp/db"
	"crudapp/routes"
	
)
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "go-crud"
)

func generateRandomKey(length int) (string, error) {
    randomBytes := make([]byte, length)
    _, err := rand.Read(randomBytes)
    if err != nil {
        return "", err
    }

    return base64.URLEncoding.EncodeToString(randomBytes), nil
}
func main() {
	// db.InitDB("user=youruser dbname=mydatabase sslmode=disable")
	// key, err := generateRandomKey(32) // 32 bytes for a 256-bit key
    // if err != nil {
    //     fmt.Println("Error generating key:", err)
    //     return
    // }

    // fmt.Println("Generated Key:", key)
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db.InitDB(connStr)

	r := routes.SetupRoutes()

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	
}
