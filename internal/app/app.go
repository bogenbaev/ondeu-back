package app

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/handler"
	remote2 "gitlab.com/a5805/ondeu/ondeu-back/internal/remote"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/repository"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/server"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/service"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/gocloak/implementation"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func Run() {
	logrus.SetReportCaller(true)

	cfg := initConfigs()
	setLogLevel(cfg.LogLevel)

	keycloak := implementation.Keycloak(cfg.Keycloak.Host, cfg.Keycloak.Realm)

	db := repository.NewPostgresRepository(repository.Config{
		Host:     cfg.Database.Host,
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
		Dbname:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})

	objectStorageConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			cfg.ObjectStorage.ClientKey,
			cfg.ObjectStorage.ClientSecret,
			""),
		Endpoint:         aws.String(cfg.ObjectStorage.Endpoint),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false), // // Configures to use subdomain/virtual calling format. Depending on your version, alternatively use o.UsePathStyle = false
	}
	newSession, err := session.NewSession(objectStorageConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	s3Client := s3.New(newSession)

	repo := repository.NewRepository(db)
	remote := remote2.NewRemote(s3Client, cfg.ObjectStorage)
	services := service.NewServices(cfg, keycloak, repo, remote)
	handlers := handler.NewHandler(services, repo, keycloak)
	srv := new(server.Server)

	go func() {
		if err = srv.Run(cfg.Port, handlers.Init()); err != nil {
			logrus.Errorf("error occured while running http server %s/n", err.Error())
		}
	}()

	logrus.Printf("server is starting at port: %s", cfg.Port)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	logrus.Printf("server is stopping at port: %s", cfg.Port)

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}

func initConfigs() *modules.AppConfigs {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Error("Error loading .env file")
	}

	keycloak := &modules.Keycloak{
		Host:              os.Getenv("KEYCLOAK_HOST"),
		ClientID:          os.Getenv("KEYCLOAK_CLIENT_ID"),
		Realm:             os.Getenv("KEYCLOAK_REALM"),
		AdminClientID:     os.Getenv("KEYCLOAK_ADMIN_CLIENT_ID"),
		AdminClientSecret: os.Getenv("KEYCLOAK_ADMIN_CLIENT_SECRET"),
	}

	databasePort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		databasePort = 5432
		fmt.Printf("DB_PORT is not set, using default port: %d\n", databasePort)
	}
	database := &modules.Postgre{
		Host:     os.Getenv("DB_HOST"),
		Port:     databasePort,
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	objectStorage := &modules.ObjectStorage{
		Endpoint:     os.Getenv("SPACES_ENDPOINT"),
		Bucket:       os.Getenv("SPACES_BUCKET"),
		ClientName:   os.Getenv("SPACES_CLIENT_NAME"),
		ClientSecret: os.Getenv("SPACES_CLIENT_SECRET"),
		ClientKey:    os.Getenv("SPACES_CLIENT_KEY"),
	}

	return &modules.AppConfigs{
		Port:          os.Getenv("PORT"),
		LogLevel:      os.Getenv("LOG_LEVEL"),
		Keycloak:      keycloak,
		Database:      database,
		ObjectStorage: objectStorage,
	}
}

func setLogLevel(level string) {
	switch level {
	case "debug":
		{
			logrus.SetLevel(logrus.DebugLevel)
			break
		}
	case "info":
		{
			logrus.SetLevel(logrus.InfoLevel)
			break
		}
	case "warn":
		{
			logrus.SetLevel(logrus.WarnLevel)
			break
		}
	case "error":
		{
			logrus.SetLevel(logrus.ErrorLevel)
			break
		}
	case "fatal":
		{
			logrus.SetLevel(logrus.FatalLevel)
			break
		}
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}
