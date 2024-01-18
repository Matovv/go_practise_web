package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Matovv/go_practise_web/mvc_mongo/controllers"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDbClient *mongo.Client

const port int = 8080

func main() {
	connectMongoDB()
	fmt.Println("\nListening on port", port)

	router := httprouter.New()
	userController := controllers.NewUserController(mongoDbClient)
	router.GET("/user/:id", userController.GetUser)
	router.POST("/user", userController.CreateUser)
	//router.DELETE("/user/:id", createUserHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))

}

func connectMongoDB() {
	fmt.Println("Connecting to MongoDB ...")

	uri := "mongodb://localhost:27017"
	// Set up the MongoDB client options.
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the MongoDB server.
	mongoDbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server to check the connection.
	err = mongoDbClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	/*
	// Close the connection when done.
	err = mongoDbClient.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	*/
}
