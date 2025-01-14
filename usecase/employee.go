package usecase

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TheMikeKaisen/Go_MongoDB/model"
	"github.com/TheMikeKaisen/Go_MongoDB/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeService struct {
	MongoCollection *mongo.Collection
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (svc *EmployeeService) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	//set response header
	w.Header().Set("Content-Type", "application/json")

	// define response
	var res Response
	defer json.NewEncoder(w).Encode(&res)	// will run in the end when all function have been executed

	// define employee model to store the request body
	var emp model.Employee
	err := json.NewDecoder(r.Body).Decode(&emp) // decode request body and store it in emp model
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("Error while decoding the request.")
	}

	// create an employee id
	emp.EmployeeId = uuid.NewString()

	// create a collection
	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	
	// insert employee 
	insertId, err := repo.InsertEmployee(&emp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("insert error: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = emp.EmployeeId
	w.WriteHeader(http.StatusOK)

	log.Println("employee inserted with id ", insertId, emp)


}


func (svc *EmployeeService) GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var res Response
	defer json.NewEncoder(w).Encode(&res)

	params := mux.Vars(r)
	empId := params["id"]
	log.Println("employee id: ", empId)

	if empId == ""{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid employee id")
		res.Error = "invalid employee id"
		return
	}

	
	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection} 

	emp, err := repo.FindEmployeeById(empId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get employee error: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = emp
	w.WriteHeader(http.StatusOK)	
}


// GetAllEmployee returns all employees in the database.
func (svc *EmployeeService) GetAllEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// response
	var res Response
	defer json.NewEncoder(w).Encode(&res)

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}

	emp, err := repo.FindAllEmployee()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error while getting all employee: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = emp
	w.WriteHeader(http.StatusOK)
}


func (svc *EmployeeService) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// response
	var res Response
	defer json.NewEncoder(w).Encode(&res)

	// params
	params := mux.Vars(r)
	empId := params["id"]

	if empId == ""{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid employee id")
		res.Error = "invalid employee id"
		return
	}

	var emp model.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("Error while decoding request body")
	}
	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	count, err := repo.UpdateEmployeeById(empId, &emp)
	if err !=nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error while updating employee: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = count
	w.WriteHeader(http.StatusOK)
}
func (svc *EmployeeService) DeleteEmployeeById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// response
	var res Response
	defer json.NewEncoder(w).Encode(&res)

	// extract id
	params := mux.Vars(r)
	empId := params["id"]

	// return error if empId is empty
	if empId == ""{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid employee id")
		res.Error = "invalid employee id"
		return
	}

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}

	// delete the emp
	count, err := repo.DeleteEmployeeById(empId)
	if err !=nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error while deleting employee: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = count
	w.WriteHeader(http.StatusOK)

}
func (svc *EmployeeService) DeleteAllEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	//response
	var res Response
	defer json.NewEncoder(w).Encode(&res)

	// create collection
	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	
	//delete employees
	count, err := repo.DeleteAllEmployees()
	if err !=nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error while deleting all employees: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = count
	w.WriteHeader(http.StatusOK)
}