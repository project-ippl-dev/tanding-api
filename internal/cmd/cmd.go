package cmd

import (
	"database/sql"
	"fmt"
	"github.com/project-ippl-dev/tanding-api/internal/document"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/project-ippl-dev/tanding-api/config"
	"github.com/project-ippl-dev/tanding-api/internal/auth"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/file"
	middlewareApp "github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/user"
)

func Run() error {
	dbConn, err := config.DatabaseConnection()
	if err != nil {
		return err
	}
	defer dbConn.Close()
	rdb, err := config.RedisConnection()
	if err != nil {
		return err
	}
	defer rdb.Close()
	sess := config.S3Connection()
	e := echo.New()

	routing(dbConn, rdb, sess, e)

	return e.Start(fmt.Sprintf("%s:%s", config.Configuration().Server.Host, config.Configuration().Server.Port))
}

func routing(dbConn *sql.DB, rdb *redis.Client, sess *session.Session, e *echo.Echo) {
	//Setup Cors
	e.Use(middleware.CORS())

	//Recover from panic
	e.Use(middleware.Recover())

	//Set Logger for every http request
	e.Use(middleware.Logger())

	//Register Static Files
	e.Static("/icon", "public/icon")

	repository := db.New(dbConn)

	//Init Middleware
	m := middlewareApp.InitMiddleware(repository)

	middlewareArgs := middlewareApp.Params{Middleware: m}

	//Init new group route
	profileRoute := e.Group("/profile", middlewareArgs.JWTMiddleware())

	//Declare Raw Repository
	userRepository := user.NewRepository(dbConn)

	//Declare Usecase
	authUsecase := auth.NewUsecase(repository, dbConn, rdb)
	userUsecase := user.NewUsecase(repository, userRepository)
	documentUsecase := document.NewUsecase(repository)
	fileUsecase := file.NewUsecase(sess)

	//Declare Handlers
	auth.RegisterHandler(authUsecase, e)
	user.RegisterHandler(userUsecase, middlewareArgs, e)
	document.RegisterHandler(documentUsecase, profileRoute)
	file.RegisterHandler(fileUsecase, middlewareArgs, e)
}
