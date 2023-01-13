package services

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/walkmanrd/assessment/models"
	"github.com/walkmanrd/assessment/repositories"
	"github.com/walkmanrd/assessment/types"
)

// ExpenseService is a struct for expense service
type ExpenseService struct {
	expenseRepository repositories.ExpenseRepository
}

// NewExpenseService is a function to create new expense service
func NewExpenseService(expenseRepository repositories.ExpenseRepository) *ExpenseService {
	return &ExpenseService{
		expenseRepository: expenseRepository,
	}
}

// GetById is a service function to get an expense by id
func (c *ExpenseService) Gets() ([]models.Expense, error) {
	expenses, err := c.expenseRepository.FindAll()

	if err != nil {
		return nil, err
	}

	return expenses, nil
}

// GetById is a service function to get an expense by id
func (c *ExpenseService) GetById(id string) (models.Expense, int, error) {
	expense, err := c.expenseRepository.FindOne(id)

	switch err {
	case sql.ErrNoRows:
		return models.Expense{}, http.StatusNotFound, errors.New("expense not found")
	case nil:
		return expense, 0, nil
	default:
		return models.Expense{}, http.StatusInternalServerError, err
	}
}

// Create is a service function to create a new expense
func (c *ExpenseService) Create(expenseRequest types.ExpenseRequest) (models.Expense, error) {
	expense, err := c.expenseRepository.Create(expenseRequest)

	if err != nil {
		return models.Expense{}, err
	}

	return expense, nil
}

// UpdateById is a service function to update an expense by id
func (c *ExpenseService) UpdateById(id string, expenseRequest types.ExpenseRequest) (models.Expense, error) {
	expense, err := c.expenseRepository.Update(id, expenseRequest)

	if err != nil {
		return models.Expense{}, err
	}

	return expense, nil
}
