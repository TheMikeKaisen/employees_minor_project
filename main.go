package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/TheMikeKaisen/Go_MongoDB/usecase"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

func init() {
	// load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load env", err)
	}

	log.Println("env loaded successfully")

	// create a mongo client
	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("error connecting to database", err)
	}

	// ping database
	err = mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Ping failed", err)
	}

	log.Println("connected to database successfully!")
}

func main() {
	// disconnect the db connection
	defer mongoClient.Disconnect(context.Background())

	// creating collection
	coll := mongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))

	// create employee service
	empService := usecase.EmployeeService{MongoCollection: coll}

	// creating Router
	r := mux.NewRouter()

	// testing if app works
	r.HandleFunc("/health", healthHandler).Methods("GET")

	r.HandleFunc("/employee", empService.GetAllEmployee).Methods("GET")
	r.HandleFunc("/employee/{id}", empService.GetEmployeeById).Methods("GET")
	r.HandleFunc("/employee", empService.CreateEmployee).Methods("POST")
	r.HandleFunc("/employee/{id}", empService.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/employee/{id}", empService.DeleteEmployeeById).Methods("DELETE")
	r.HandleFunc("/employee", empService.DeleteAllEmployee).Methods("DELETE")

	// listening on port 8000
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal("error connecting to the server")
	}
	log.Println("Running on port 8000")
}


func healthHandler(w http.ResponseWriter, r *http.Request) {
// of 200 (OK) and a simple text message "Running..." to indicate that the API
// is up and running.

	w.Write([]byte("Running..."))
}
