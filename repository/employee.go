package repository

import (
	"context"
	"fmt"

	"github.com/TheMikeKaisen/Go_MongoDB/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeRepo struct {
	MongoCollection *mongo.Collection // represents the memory of a collection in mongodb
}

func (r *EmployeeRepo) InsertEmployee(emp *model.Employee) (interface{}, error) {

	// convert emp into bson format and insert it into the collection
	result, err := r.MongoCollection.InsertOne(context.Background(), emp)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *EmployeeRepo) FindEmployeeById(empId string) (*model.Employee, error) {
	var emp model.Employee

	err := r.MongoCollection.FindOne(context.Background(),
		bson.D{{Key: "employee_id", Value: empId}}).Decode(&emp)

	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepo) FindAllEmployee() ([]model.Employee, error) {

	// find all employees from database and put them into result(the result is a mongo.cursor-> it needs to be decoded)
	result, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var emps []model.Employee

	// decoding mongo.cursor into emps
	err = result.All(context.Background(), &emps)
	if err != nil {
		return nil, fmt.Errorf("results decode error %s", err.Error())
	}
	return emps, nil

}

func (r *EmployeeRepo) UpdateEmployeeById(empId string, updateEmp *model.Employee) (int64, error) {

	result, err := r.MongoCollection.UpdateOne(
		context.Background(),
		bson.D{{Key: "employee_id", Value: empId}},
		bson.D{{Key: "$set", Value: updateEmp}})

	if err != nil {
		return 0, nil
	}

	return result.ModifiedCount, nil
}


func (r *EmployeeRepo) DeleteEmployeeById(empId string) (int64, error) {

	result, err := r.MongoCollection.DeleteOne(
		context.Background(),
		bson.D{{Key: "employee_id", Value: empId}})

	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *EmployeeRepo) DeleteAllEmployees() (int64, error) {
	result, err := r.MongoCollection.DeleteMany(
		context.Background(),
		bson.D{})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
