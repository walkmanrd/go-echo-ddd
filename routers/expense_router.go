package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/walkmanrd/assessment/controllers"
)

// ExpenseRouter is a function to set expense routes on resource path /expenses
func ExpenseRouter(e *echo.Group) {

	// ExpenseController is a struct for expense controller
	var expenseController controllers.ExpenseController

	// Setting up routes
	e.GET("", expenseController.Index)
	e.GET("/:id", expenseController.Show)
	e.POST("", expenseController.Store)
	e.PUT("/:id", expenseController.Update)
}
