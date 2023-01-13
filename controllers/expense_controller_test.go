//go:build unit

package controllers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/walkmanrd/assessment/repositories"
	"github.com/walkmanrd/assessment/services"
	"github.com/walkmanrd/assessment/validators"
)

var requestBody = `{"id":"1","title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}`

func setupTest(db *sql.DB) *ExpenseController {
	expenseRepository := repositories.NewExpenseRepository(db)
	expenseService := services.NewExpenseService(*expenseRepository)

	expenseController := &ExpenseController{
		expenseService: *expenseService,
	}

	return expenseController
}

func TestCreateExpense(t *testing.T) {
	e := echo.New()
	e.Validator = &validators.CustomValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow("1", "strawberry smoothie", 79.0, "night market promotion discount 10 bath", `{"food","beverage"}`)
	db, mock, err := sqlmock.New()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO expenses (id, title, amount, note, tags) values (DEFAULT, $1, $2, $3, $4)
	RETURNING id, title, amount, note, tags`)).
		WithArgs("strawberry smoothie", 79.0, "night market promotion discount 10 bath", `{"food","beverage"}`).
		WillReturnRows(mockRows)

	if err != nil {
		t.Fatalf("Error '%s' was not expected when opening a stub database connection", err)
	}

	expenseController := setupTest(db)

	if assert.NoError(t, expenseController.Store(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, requestBody, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetExpenseById(t *testing.T) {
	e := echo.New()
	e.Validator = &validators.CustomValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodGet, "/expenses/1", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow("1", "strawberry smoothie get by id", 99.0, "night market promotion discount 10 bath get by id", `{"food","beverage","get by id"}`)

	db, mock, err := sqlmock.New()
	mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, title, amount, note, tags FROM expenses WHERE id = $1`)).
		ExpectQuery().
		WithArgs("1").
		WillReturnRows(mockRows)

	if err != nil {
		t.Fatalf("Error '%s' was not expected when opening a stub database connection", err)
	}

	expenseRepository := repositories.NewExpenseRepository(db)
	expenseService := services.NewExpenseService(*expenseRepository)

	expenseController := &ExpenseController{
		expenseService: *expenseService,
	}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	var expectTestGetById = `{"id":"1","title":"strawberry smoothie get by id","amount":99,"note":"night market promotion discount 10 bath get by id","tags":["food","beverage","get by id"]}`

	if assert.NoError(t, expenseController.Show(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectTestGetById, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUpdateExpenseById(t *testing.T) {
	e := echo.New()
	e.Validator = &validators.CustomValidator{Validator: validator.New()}
	var updateBody = `{"id":"1","title":"strawberry smoothie update","amount":100,"note":"night market promotion discount 10 bath update","tags":["food","beverage","update"]}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updateBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "strawberry smoothie update", 100, "night market promotion discount 10 bath update", `{"food","beverage","update"}`)
	db, mock, err := sqlmock.New()

	mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1 RETURNING id, title, amount, note, tags;`)).
		ExpectQuery().
		WithArgs("1", "strawberry smoothie update", 100.0, "night market promotion discount 10 bath update", `{"food","beverage","update"}`).
		WillReturnRows(mockRows)

	if err != nil {
		t.Fatalf("Error '%s' was not expected when opening a stub database connection", err)
	}
	expenseController := setupTest(db)

	if assert.NoError(t, expenseController.Update(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, updateBody, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetExpenses(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "strawberry smoothie 1", 71, "night market promotion discount 10 bath 1", `{"food","beverage","1"}`).
		AddRow("2", "strawberry smoothie 2", 72, "night market promotion discount 10 bath 2", `{"food","beverage","2"}`)

	db, mock, err := sqlmock.New()
	mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, title, amount, note, tags FROM expenses ORDER BY id ASC`)).
		ExpectQuery().
		WillReturnRows(mockRows)
	if err != nil {
		t.Fatalf("Error '%s' was not expected when opening a stub database connection", err)
	}

	var expectTestGetAll = `[{"id":"1","title":"strawberry smoothie 1","amount":71,"note":"night market promotion discount 10 bath 1","tags":["food","beverage","1"]},{"id":"2","title":"strawberry smoothie 2","amount":72,"note":"night market promotion discount 10 bath 2","tags":["food","beverage","2"]}]`

	expenseController := setupTest(db)
	c := e.NewContext(req, rec)
	c.SetPath("/expenses")

	if assert.NoError(t, expenseController.Index(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectTestGetAll, strings.TrimSpace(rec.Body.String()))
	}
}
