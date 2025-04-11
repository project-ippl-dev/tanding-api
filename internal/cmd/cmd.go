package cmd

import (
	"database/sql"
	"github.com/dytlan/tanding-api/config"
	"github.com/dytlan/tanding-api/internal/db"
	"github.com/dytlan/tanding-api/internal/test"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() error {
	dbConn, err := config.DatabaseConnection()
	if err != nil {
		return err
	}
	defer dbConn.Close()
	e := echo.New()

	routing(dbConn, e)

	return e.Start(config.Configuration().Server.Port)
}

func routing(dbConn *sql.DB, e *echo.Echo) {
	//Setup Cors
	e.Use(middleware.CORS())

	//Recover from panic
	e.Use(middleware.Recover())

	//Set Logger for every http request
	e.Use(middleware.Logger())

	repository := db.New(dbConn)

	//Declare Usecase
	testUsecase := test.NewUsecase(repository)

	//Declare Handlers
	test.RegisterHandler(testUsecase, e)
}
