package types

// ExpenseRequest is a type for expense request
type ExpenseRequest struct {
	Title  string   `json:"title" validate:"required"`
	Amount float64  `json:"amount" validate:"required,number"`
	Note   string   `json:"note" validate:"required"`
	Tags   []string `json:"tags" validate:"required,min=1"`
}
