package model

// Employee represents an employee record in MongoDB.
type Employee struct {
	EmployeeId   string `json:"employee_id,omitempty" bson:"employee_id,omitempty"`
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Department   string `json:"department,omitempty" bson:"department,omitempty"`

	// New fields
	MobileNumber string `json:"mobile_number,omitempty" bson:"mobile_number,omitempty"`  // e.g., "+1-555-1234"
	Gender       string `json:"gender,omitempty" bson:"gender,omitempty"`          // e.g., "Male", "Female", "Non-binary"
	Email        string `json:"email,omitempty" bson:"email,omitempty"`            // e.g., "alice@example.com"
	Age          int    `json:"age,omitempty" bson:"age,omitempty"`              // Employee's age in years
}