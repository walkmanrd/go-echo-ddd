package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	// "github.com/go-errors/errors"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/walkmanrd/assessment/configs"
	"github.com/walkmanrd/assessment/routers"
	"github.com/walkmanrd/assessment/types"
	"github.com/walkmanrd/assessment/validators"

	_ "github.com/lib/pq"
)

// init is a function that run before main
func init() {
	db := configs.ConnectDatabase()
	configs.AutoMigrate(db)
	defer db.Close()
}

// AuthHeader is a middleware to check authorization header
func AuthHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")
		authTokenCheck := os.Getenv("AUTH_TOKEN")

		if authorization == authTokenCheck {
			return next(c)
		}
		return c.JSON(http.StatusUnauthorized, types.Error{Message: "Unauthorized"})
	}
}

// main is a function that run after init
func main() {
	// Echo instance
	e := echo.New()
	e.Validator = &validators.CustomValidator{Validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes Public
	routers.HealthCheckRouter(e)

	// Routes Private
	g := e.Group("/expenses")
	g.Use(AuthHeader)
	routers.ExpenseRouter(g)

	// Start server
	port := os.Getenv("PORT")

	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("shutting complete bye bye")
}
