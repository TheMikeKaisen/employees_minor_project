package usecase

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/TheMikeKaisen/Go_MongoDB/model"
	"github.com/TheMikeKaisen/Go_MongoDB/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// EmployeeService provides employee-related operations.
type EmployeeService struct {
	MongoCollection *mongo.Collection
}

// Response is the standard API response.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// validateEmployee applies basic business validations.
func validateEmployee(emp *model.Employee) error {
	if emp.Name == "" {
		return errors.New("name is required")
	}
	if emp.Department == "" {
		return errors.New("department is required")
	}
	// Age must be positive
	if emp.Age <= 0 {
		return errors.New("age must be a positive integer")
	}
	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if emp.Email == "" || !emailRegex.MatchString(emp.Email) {
		return errors.New("invalid email address")
	}
	// Mobile number length check (minimum 7 digits)
	mobileRegex := regexp.MustCompile(`^\+?[0-9]{7,15}$`)
	if emp.MobileNumber == "" || !mobileRegex.MatchString(emp.MobileNumber) {
		return errors.New("invalid mobile number")
	}
	// Gender limited values
	switch emp.Gender {
	case "Male", "Female", "Non-binary", "Other":
		// ok
	default:
		return errors.New("gender must be Male, Female, Non-binary, or Other")
	}
	return nil
}

// CreateEmployee handles POST /employees
func (svc *EmployeeService) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	defer json.NewEncoder(w).Encode(&res)

	var emp model.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = "invalid request payload"
		return
	}

	// assign new ID and validate
	emp.EmployeeId = uuid.NewString()
	if err := validateEmployee(&emp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = err.Error()
		return
	}

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	insertID, err := repo.InsertEmployee(&emp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("insert error:", err)
		res.Error = err.Error()
		return
	}

	res.Data = map[string]interface{}{"inserted_id": insertID, "employee_id": emp.EmployeeId}
	w.WriteHeader(http.StatusCreated)
}

// GetEmployeeById handles GET /employees/{id}
func (svc *EmployeeService) GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	defer json.NewEncoder(w).Encode(&res)

	empId := mux.Vars(r)["id"]
	if empId == "" {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = "employee ID is required"
		return
	}

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	emp, err := repo.FindEmployeeById(empId)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, mongo.ErrNoDocuments) {
			status = http.StatusNotFound
		}
		w.WriteHeader(status)
		res.Error = err.Error()
		return
	}

	res.Data = emp
	w.WriteHeader(http.StatusOK)
}

// GetAllEmployee handles GET /employees
func (svc *EmployeeService) GetAllEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	defer json.NewEncoder(w).Encode(&res)

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	temps, err := repo.FindAllEmployee()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = err.Error()
		return
	}

	res.Data = temps
	w.WriteHeader(http.StatusOK)
}

// UpdateEmployee handles PUT /employees/{id}
func (svc *EmployeeService) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	defer json.NewEncoder(w).Encode(&res)

	empId := mux.Vars(r)["id"]
	if empId == "" {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = "employee ID is required"
		return
	}

	var emp model.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = "invalid request payload"
		return
	}

	if err := validateEmployee(&emp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = err.Error()
		return
	}

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	count, err := repo.UpdateEmployeeById(empId, &emp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = err.Error()
		return
	}

	res.Data = map[string]int64{"modified_count": count}
	w.WriteHeader(http.StatusOK)
}

// DeleteEmployeeById handles DELETE /employees/{id}
func (svc *EmployeeService) DeleteEmployeeById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	defer json.NewEncoder(w).Encode(&res)

	empId := mux.Vars(r)["id"]
	if empId == "" {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = "employee ID is required"
		return
	}

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	count, err := repo.DeleteEmployeeById(empId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = err.Error()
		return
	}

	res.Data = map[string]int64{"deleted_count": count}
	w.WriteHeader(http.StatusOK)
}

// DeleteAllEmployee handles DELETE /employees
func (svc *EmployeeService) DeleteAllEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response
	defer json.NewEncoder(w).Encode(&res)

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	count, err := repo.DeleteAllEmployees()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = err.Error()
		return
	}

	res.Data = map[string]int64{"deleted_count": count}
	w.WriteHeader(http.StatusOK)
}
