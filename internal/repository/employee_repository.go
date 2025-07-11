package repository

import "database/sql"

type Employee struct {
	NPK      string `json:"npk"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

type EmployeeRepository struct {
	DB *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

func (r *EmployeeRepository) GetAllEmployees() ([]Employee, error) {
	rows, err := r.DB.Query("SELECT npk, name, position FROM employee WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.NPK, &e.Name, &e.Position); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}