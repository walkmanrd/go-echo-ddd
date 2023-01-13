package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/walkmanrd/assessment/services"
	"github.com/walkmanrd/assessment/types"
)

// ExpenseController is a struct for expense controller
type ExpenseController struct {
	expenseRequest types.ExpenseRequest
	expenseService services.ExpenseService
}

// bindAndValidateRequest is a function to bind and validate request
func bindAndValidateRequest(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	return nil
}

// GET /expenses
// Index is a function to get all expenses
func (c *ExpenseController) Index(e echo.Context) error {
	expenses, err := c.expenseService.Gets()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, types.Error{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, expenses)
}

// GET /expenses/:id
// Show is a function to get an expense by id
func (c *ExpenseController) Show(e echo.Context) error {
	id := e.Param("id")

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return e.JSON(http.StatusBadRequest, types.Error{Message: "invalid parameter id"})
	}

	expense, status, err := c.expenseService.GetById(id)

	if err != nil {
		return e.JSON(status, types.Error{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, expense)
}

// POST /expenses
// Store is a function to create a new expense
func (c *ExpenseController) Store(e echo.Context) error {
	if err := bindAndValidateRequest(e, &c.expenseRequest); err != nil {
		return e.JSON(http.StatusBadRequest, types.Error{Message: err.Error()})
	}

	expense, err := c.expenseService.Create(c.expenseRequest)

	if err != nil {
		return e.JSON(http.StatusInternalServerError, types.Error{Message: err.Error()})
	}

	return e.JSON(http.StatusCreated, expense)
}

// PUT /expenses/:id
// Update is a function to get an expense by id
func (c *ExpenseController) Update(e echo.Context) error {
	id := e.Param("id")

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return e.JSON(http.StatusBadRequest, types.Error{Message: "invalid parameter id"})
	}

	if err := bindAndValidateRequest(e, &c.expenseRequest); err != nil {
		return e.JSON(http.StatusBadRequest, types.Error{Message: err.Error()})
	}

	newExpense, err := c.expenseService.UpdateById(id, c.expenseRequest)

	if err != nil {
		return e.JSON(http.StatusInternalServerError, types.Error{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, newExpense)
}
