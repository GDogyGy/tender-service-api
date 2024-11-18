package service

import "TenderServiceApi/internal/model"

type EmployeeService struct {
	EmployeeRepository EmployeeRepository
}

type EmployeeRepository interface {
	GetEmployeeByID(id string) (model.Employee, error)
}

func NewEmployeeService(EmployeeRepository EmployeeRepository) *EmployeeService {
	return &EmployeeService{EmployeeRepository: EmployeeRepository}
}

func (s *EmployeeService) GetEmployeeByID(id string) (model.Employee, error) {
	return s.EmployeeRepository.GetEmployeeByID(id)
}
