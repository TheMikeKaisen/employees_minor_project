package repository

import (
	"context"
	"log"
	"testing"

	"github.com/TheMikeKaisen/Go_MongoDB/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func newMongoClient() *mongo.Client {
	mongoTestClient, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI("mongodb+srv://mike:nair12345678@cluster0.nw7sz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))

	if err != nil {
		log.Fatal("error while connecting mongodb", err)
	}

	err = mongoTestClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("ping failed", err)
	}
	log.Println("ping success")
	return mongoTestClient
}

func TestMongoOperations(t *testing.T) {

	// create a client
	mongoTestClient := newMongoClient()
	defer mongoTestClient.Disconnect(context.Background())

	// dummy data
	emp1 := uuid.New().String()
	// emp2 :=uuid.New().String()

	// connect to a collection
	coll := mongoTestClient.Database("companydb").Collection("employee_test")
	empRepo := EmployeeRepo{MongoCollection: coll}

	t.Run("Insert Employee 1", func(t *testing.T) {
		emp := model.Employee{
			EmployeeId: emp1,
			Name:       "Tony Stark",
			Department: "physics",
		}

		result, err := empRepo.InsertEmployee(&emp)
		if err != nil {
			t.Fatal("Error inserting employee 1")
		}
		t.Log("Employee was successfully inserted!", result)
	})

	t.Run("Get Employee 1", func(t *testing.T) {
		result, err := empRepo.FindEmployeeById(emp1)
		if err != nil {
			t.Fatal("Error retrieving employee 1")
		}
		t.Log("Employee was successfully retrieved!", result.Name)
	})

	t.Run("Get All Employees", func(t *testing.T) {
		result, err := empRepo.FindAllEmployee()
		if err != nil {
			t.Fatal("Error while receiving employees!")
		}
		t.Log("List of Employees: ", result)
	})
	t.Run("Get All Employees", func(t *testing.T) {
		newEmp := model.Employee{
			EmployeeId: emp1,
			Name:       "Steven Rogers",
			Department: "physics",
		}
		result, err := empRepo.UpdateEmployeeById(emp1, &newEmp)
		if err != nil {
			t.Fatal("Error while updating employees!")
		}
		t.Log("Updated Employee Count: ", result)
	})

	t.Run("Delete employee by id", func(t *testing.T) {
		result, err := empRepo.DeleteEmployeeById(emp1)
		if err != nil {
			t.Fatal("Error while deleting the employee!")
		}
		t.Log("Employee Deleted Count: ", result)
	})

	t.Run("Delete employee by id", func(t *testing.T) {
		result, err := empRepo.DeleteAllEmployees()
		if err != nil {
			t.Fatal("Error while deleting the employees!")
		}
		t.Log("Employees Deleted Count: ", result)
	})

}
