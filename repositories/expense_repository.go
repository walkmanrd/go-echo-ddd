package repositories

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/walkmanrd/assessment/configs"
	"github.com/walkmanrd/assessment/models"
	"github.com/walkmanrd/assessment/types"
)

// ExpenseRepository is a repository for expense
type ExpenseRepository struct {
	db *sql.DB
}

// NewExpenseRepository is a function to create new expense repository
func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{
		db: db,
	}
}

// FindAll is a function to get all expenses
func (r *ExpenseRepository) FindAll() ([]models.Expense, error) {
	if r.db == nil {
		r.db = configs.ConnectDatabase()
	}

	stmt, err := r.db.Prepare("SELECT id, title, amount, note, tags FROM expenses ORDER BY id ASC")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	expenses := []models.Expense{}

	for rows.Next() {
		expense := models.Expense{}
		err := rows.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, &expense.Tags)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

// FindOne is a function to get an expenses by id
func (r *ExpenseRepository) FindOne(id string) (models.Expense, error) {
	if r.db == nil {
		r.db = configs.ConnectDatabase()
	}

	stmt, err := r.db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return models.Expense{}, err
	}

	row := stmt.QueryRow(id)
	expense := models.Expense{}

	err = row.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, &expense.Tags)
	if err != nil {
		return models.Expense{}, err
	}

	return expense, nil
}

// Create is a function to create a new expense
func (r *ExpenseRepository) Create(expenseRequest types.ExpenseRequest) (models.Expense, error) {
	if r.db == nil {
		r.db = configs.ConnectDatabase()
	}

	sqlCommand := `
	INSERT INTO expenses (id, title, amount, note, tags) values (DEFAULT, $1, $2, $3, $4)
	RETURNING id, title, amount, note, tags
	`
	expense := models.Expense{
		Title:  expenseRequest.Title,
		Amount: expenseRequest.Amount,
		Note:   expenseRequest.Note,
		Tags:   expenseRequest.Tags,
	}

	tags := pq.Array(expense.Tags)
	row := r.db.QueryRow(sqlCommand, &expense.Title, &expense.Amount, &expense.Note, tags)
	err := row.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, &expense.Tags)

	if err != nil {
		fmt.Println("can't scan id on ExpenseRepository", err)
		return models.Expense{}, err
	}

	return expense, nil
}

// Update is a function to update an expense by id
func (r *ExpenseRepository) Update(id string, expenseRequest types.ExpenseRequest) (models.Expense, error) {
	if r.db == nil {
		r.db = configs.ConnectDatabase()
	}

	sqlCommand := `UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1 RETURNING id, title, amount, note, tags;`

	expense := models.Expense{
		ID:     id,
		Title:  expenseRequest.Title,
		Amount: expenseRequest.Amount,
		Note:   expenseRequest.Note,
		Tags:   expenseRequest.Tags,
	}

	stmt, err := r.db.Prepare(sqlCommand)

	if err != nil {
		fmt.Println("can't prepare statement on ExpenseRepository", err)
		return models.Expense{}, err
	}

	row := stmt.QueryRow(id, &expense.Title, &expense.Amount, &expense.Note, &expense.Tags)
	err = row.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, &expense.Tags)

	if err != nil {
		fmt.Println("can't scan id on ExpenseRepository", err)
		return models.Expense{}, err
	}

	return expense, nil
}
