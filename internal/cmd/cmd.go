package cmd

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/project-ippl-dev/tanding-api/config"
	"github.com/project-ippl-dev/tanding-api/internal/accomplishment"
	"github.com/project-ippl-dev/tanding-api/internal/auth"
	"github.com/project-ippl-dev/tanding-api/internal/bracket"
	"github.com/project-ippl-dev/tanding-api/internal/certificate"
	"github.com/project-ippl-dev/tanding-api/internal/club"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/document"
	"github.com/project-ippl-dev/tanding-api/internal/event"
	"github.com/project-ippl-dev/tanding-api/internal/eventRegistration"
	"github.com/project-ippl-dev/tanding-api/internal/file"
	"github.com/project-ippl-dev/tanding-api/internal/mail"
	middlewareApp "github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/sport"
	"github.com/project-ippl-dev/tanding-api/internal/storage"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	"github.com/project-ippl-dev/tanding-api/internal/user"
)

func Run() error {
	conf := config.NewConfig()

	postgresDB, err := config.NewDatabase(conf.Database)
	if err != nil {
		return err
	}
	defer postgresDB.Close()

	rdb, err := config.NewRedisClient(conf.Redis)
	if err != nil {
		return err
	}

	defer rdb.Close()

	s3Client := config.NewS3Client(conf.S3)
	jwtClient := tools.NewJWTClient(conf.JWT)
	mailClient := mail.NewMailClient(conf.SMTP)
	storageClient := config.NewStorageClient(conf.StorageConfig)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	e := echo.New()

	routing(conf, postgresDB, rdb, s3Client, jwtClient, storageClient, mailClient, e, r)

	return e.Start(fmt.Sprintf("%s:%d", conf.ServerConfig.Host, conf.ServerConfig.Port))
}

func routing(conf config.Config, postgresDB *sql.DB, rdb *redis.Client, s3Client config.S3Client, jwtClient tools.JWTClient, storageClient config.StorageClient, mailClient mail.MailClient, e *echo.Echo, r *rand.Rand) {
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

	repository := db.New(postgresDB)

	//Init Middleware
	m := middlewareApp.InitMiddleware(repository)

	middlewareArgs := middlewareApp.Params{
		Middleware: m,
		JWTClient:  jwtClient,
	}

	//Init new group route
	profileRoute := e.Group("/profile", middlewareArgs.JWTMiddleware())

	//Declare Raw Repository
	userRepository := user.NewRepository(postgresDB)
	eventRegistrationRepository := eventRegistration.NewRepository(postgresDB)
	sportRepository := sport.NewRepository(postgresDB)
	clubRepository := club.NewRepository(postgresDB)
	eventRepository := event.NewRepository(postgresDB)
	bracketRepository := bracket.NewRepository(postgresDB)

	//Declare Usecase
	authUsecase := auth.NewUsecase(repository, postgresDB, rdb, mailClient, jwtClient)
	accomplishmentUsecase := accomplishment.NewUsecase(repository)
	eventRegistrationUsecase := eventRegistration.NewUsecase(repository, eventRegistrationRepository)
	userUsecase := user.NewUsecase(repository, userRepository)
	sportUsecase := sport.NewUsecase(repository, sportRepository)
	documentUsecase := document.NewUsecase(repository)
	fileUsecase := file.NewUsecase(s3Client)
	storageUsecase := storage.NewUsecase(storageClient)
	clubUsecase := club.NewUsecase(repository, clubRepository)
	eventUsecase := event.NewUsecase(repository, eventRepository)
	bracketUsecase := bracket.NewUsecase(repository, bracketRepository, r)
	certificateUsecase := certificate.NewUsecase(repository)

	//Declare Handlers
	auth.RegisterHandler(authUsecase, jwtClient, conf.ServerConfig, e)
	accomplishment.RegisterHandler(accomplishmentUsecase, jwtClient, profileRoute)
	bracket.RegisterHandler(bracketUsecase, middlewareArgs, e)
	user.RegisterHandler(userUsecase, middlewareArgs, jwtClient, e)
	sport.RegisterHandler(sportUsecase, middlewareArgs, e)
	document.RegisterHandler(documentUsecase, jwtClient, profileRoute)
	eventRegistration.RegisterHandler(eventRegistrationUsecase, middlewareArgs, jwtClient, e)
	file.RegisterHandler(fileUsecase, middlewareArgs, e)
	storage.RegisterHandler(storageUsecase, middlewareArgs, e)
	club.RegisterHandler(clubUsecase, middlewareArgs, jwtClient, e)
	event.RegisterHandler(eventUsecase, middlewareArgs, jwtClient, e)
	certificate.RegisterHandler(certificateUsecase, middlewareArgs, e)
}
