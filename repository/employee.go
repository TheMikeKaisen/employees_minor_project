package repository

import (
	"context"
	"fmt"

	"github.com/TheMikeKaisen/Go_MongoDB/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// EmployeeRepo handles CRUD operations for Employee documents in MongoDB.
type EmployeeRepo struct {
	MongoCollection *mongo.Collection // represents the MongoDB collection
}

// InsertEmployee inserts a new employee document and returns its generated ID.
func (r *EmployeeRepo) InsertEmployee(emp *model.Employee) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), emp)
	if err != nil {
		return nil, fmt.Errorf("failed to insert employee: %w", err)
	}
	return result.InsertedID, nil
}

// FindEmployeeById retrieves a single employee by its EmployeeId field.
func (r *EmployeeRepo) FindEmployeeById(empId string) (*model.Employee, error) {
	var emp model.Employee
	filter := bson.D{{Key: "employee_id", Value: empId}}
	err := r.MongoCollection.FindOne(context.Background(), filter).Decode(&emp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no employee found with id %s", empId)
		}
		return nil, fmt.Errorf("error finding employee: %w", err)
	}
	return &emp, nil
}

// FindAllEmployee retrieves all employee documents from the collection.
func (r *EmployeeRepo) FindAllEmployee() ([]model.Employee, error) {
	cursor, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to find employees: %w", err)
	}
	var emps []model.Employee
	if err := cursor.All(context.Background(), &emps); err != nil {
		return nil, fmt.Errorf("failed to decode employees: %w", err)
	}
	return emps, nil
}

// UpdateEmployeeById updates an existing employee document by its EmployeeId.
func (r *EmployeeRepo) UpdateEmployeeById(empId string, updateEmp *model.Employee) (int64, error) {
	filter := bson.D{{Key: "employee_id", Value: empId}}
	update := bson.D{{Key: "$set", Value: updateEmp}}

	result, err := r.MongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return 0, fmt.Errorf("failed to update employee: %w", err)
	}
	return result.ModifiedCount, nil
}

// DeleteEmployeeById removes a single employee document by its EmployeeId.
func (r *EmployeeRepo) DeleteEmployeeById(empId string) (int64, error) {
	filter := bson.D{{Key: "employee_id", Value: empId}}
	result, err := r.MongoCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete employee: %w", err)
	}
	return result.DeletedCount, nil
}

// DeleteAllEmployees removes all documents from the employee collection.
func (r *EmployeeRepo) DeleteAllEmployees() (int64, error) {
	result, err := r.MongoCollection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		return 0, fmt.Errorf("failed to delete all employees: %w", err)
	}
	return result.DeletedCount, nil
}
