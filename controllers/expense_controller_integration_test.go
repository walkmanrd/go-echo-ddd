//go:build integration

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/walkmanrd/assessment/models"
	"github.com/walkmanrd/assessment/validators"

	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

const serverPort = 2565

func SetupSuite() *echo.Echo {
	eh := echo.New()
	eh.Validator = &validators.CustomValidator{Validator: validator.New()}
	eh.Use(middleware.Logger())
	eh.Use(middleware.Recover())
	var expenseController ExpenseController

	// Setting up routes
	eh.GET("/expenses", expenseController.Index)
	eh.GET("/expenses/:id", expenseController.Show)
	eh.POST("/expenses", expenseController.Store)
	eh.PUT("/expenses/:id", expenseController.Update)

	go func(e *echo.Echo) {
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	return eh
}

func seedExpense(t *testing.T) models.Expense {
	reqBody := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.Nil(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}
	resp, err := client.Do(req)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	var expense models.Expense
	err = json.Unmarshal([]byte(resBody), &expense)
	assert.NoError(t, err)

	if assert.NotEqual(t, 0, expense.ID) {
		assert.Equal(t, "strawberry smoothie", expense.Title)
		assert.Equal(t, float64(79), float64(expense.Amount))
		assert.Equal(t, "night market promotion discount 10 bath", expense.Note)
		assert.Equal(t, pq.StringArray{"food", "beverage"}, expense.Tags)
	}

	return expense
}

func getExpenseById(t *testing.T, id string) models.Expense {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/"+id, serverPort), nil)
	assert.NoError(t, err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	resp, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	var expense models.Expense
	err = json.Unmarshal(resBody, &expense)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEqual(t, 0, expense.ID)
		assert.Equal(t, "strawberry smoothie", expense.Title)
		assert.Equal(t, float64(79), float64(expense.Amount))
		assert.Equal(t, "night market promotion discount 10 bath", expense.Note)
		assert.Equal(t, pq.StringArray{"food", "beverage"}, expense.Tags)
	}

	return expense
}

func TestCreateExpense(t *testing.T) {
	// setup server
	eh := SetupSuite()

	// create expense and assertion
	seedExpense(t)

	// shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	eh.Shutdown(ctx)
}

func TestGetExpenseByID(t *testing.T) {
	// setup server
	eh := SetupSuite()

	// create expense and assertion
	newExpense := seedExpense(t)

	// get expense by id and assertion
	getExpenseById(t, newExpense.ID)

	// shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	eh.Shutdown(ctx)
}

func TestGetAllExpense(t *testing.T) {
	// setup server
	eh := SetupSuite()

	// setup seed
	seedExpense(t)

	// arrange
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses", serverPort), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err2 := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var es []models.Expense
	json.Unmarshal(resBody, &es)

	// assertions
	if assert.NoError(t, err2) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEqual(t, 0, len(es))
		assert.NotEqual(t, 0, es[0].ID)
		assert.Equal(t, "strawberry smoothie", es[0].Title)
		assert.Equal(t, float64(79), float64(es[0].Amount))
		assert.Equal(t, "night market promotion discount 10 bath", es[0].Note)
		assert.Equal(t, pq.StringArray{"food", "beverage"}, es[0].Tags)
	}

	// shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	eh.Shutdown(ctx)
}

func TestUpdateUserByID(t *testing.T) {
	// Setup server
	eh := SetupSuite()

	// create expense and assertion
	newExpense := seedExpense(t)

	// Update last Expense
	reqBody := `{
		"title": "strawberry smoothie for update",
		"amount": 80,
		"note": "night market promotion discount 10 bath for update",
		"tags": ["food", "beverage", "for update"]
	}`

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/expenses/"+newExpense.ID, serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}
	respUpdate, err := client.Do(req)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, respUpdate.StatusCode)
	}

	resBody, err := ioutil.ReadAll(respUpdate.Body)
	assert.NoError(t, err)
	respUpdate.Body.Close()

	var lastExpense models.Expense
	err = json.Unmarshal([]byte(resBody), &lastExpense)
	assert.NoError(t, err)

	assert.NotEqual(t, 0, lastExpense.ID)
	assert.Equal(t, "strawberry smoothie for update", lastExpense.Title)
	assert.Equal(t, float64(80), float64(lastExpense.Amount))
	assert.Equal(t, "night market promotion discount 10 bath for update", lastExpense.Note)
	assert.Equal(t, pq.StringArray{"food", "beverage", "for update"}, lastExpense.Tags)

	// shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	eh.Shutdown(ctx)
}
