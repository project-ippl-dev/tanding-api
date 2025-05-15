package cmd

import (
	"database/sql"
	"fmt"

	"net/http"

	"github.com/project-ippl-dev/tanding-api/internal/club"
	"github.com/project-ippl-dev/tanding-api/internal/eventRegistration"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/project-ippl-dev/tanding-api/config"
	"github.com/project-ippl-dev/tanding-api/internal/accomplishment"
	"github.com/project-ippl-dev/tanding-api/internal/auth"
	"github.com/project-ippl-dev/tanding-api/internal/bracket"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/document"
	"github.com/project-ippl-dev/tanding-api/internal/event"
	"github.com/project-ippl-dev/tanding-api/internal/file"
	middlewareApp "github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/sport"
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

	e.GET("/healthcheck", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "HEALTHY"})
	})

	repository := db.New(dbConn)

	//Init Middleware
	m := middlewareApp.InitMiddleware(repository)

	middlewareArgs := middlewareApp.Params{Middleware: m}

	//Init new group route
	profileRoute := e.Group("/profile", middlewareArgs.JWTMiddleware())

	//Declare Raw Repository
	userRepository := user.NewRepository(dbConn)
	eventRegistrationRepository := eventRegistration.NewRepository(dbConn)
	sportRepository := sport.NewRepository(dbConn)
	clubRepository := club.NewRepository(dbConn)
	eventRepository := event.NewRepository(dbConn)
	bracketRepository := bracket.NewRepository(dbConn)

	//Declare Usecase
	authUsecase := auth.NewUsecase(repository, dbConn, rdb)
	accomplishmentUsecase := accomplishment.NewUsecase(repository)
	eventRegistrationUsecase := eventRegistration.NewUsecase(repository, eventRegistrationRepository)
	userUsecase := user.NewUsecase(repository, userRepository)
	sportUsecase := sport.NewUsecase(repository, sportRepository)
	documentUsecase := document.NewUsecase(repository)
	fileUsecase := file.NewUsecase(sess)
	clubUsecase := club.NewUsecase(repository, clubRepository)
	eventUsecase := event.NewUsecase(repository, eventRepository)
	bracketUsecase := bracket.NewUsecase(repository, bracketRepository)

	//Declare Handlers
	auth.RegisterHandler(authUsecase, e)
	accomplishment.RegisterHandler(accomplishmentUsecase, profileRoute)
	bracket.RegisterHandler(bracketUsecase, middlewareArgs, e)
	user.RegisterHandler(userUsecase, middlewareArgs, e)
	sport.RegisterHandler(sportUsecase, middlewareArgs, e)
	document.RegisterHandler(documentUsecase, profileRoute)
	eventRegistration.RegisterHandler(eventRegistrationUsecase, middlewareArgs, e)
	file.RegisterHandler(fileUsecase, middlewareArgs, e)
	club.RegisterHandler(clubUsecase, middlewareArgs, e)
	event.RegisterHandler(eventUsecase, middlewareArgs, e)
}
