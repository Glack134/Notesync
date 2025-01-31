package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/polyk005/notesync"
	"github.com/polyk005/notesync/pkg/handler"
	"github.com/polyk005/notesync/pkg/repository"
	"github.com/polyk005/notesync/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(notesync.Server)
	certFile := "./certs/notesync.crt"
	keyFile := "./certs/notesync_private_key.pem"
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Fatalf("Сертификат не найден: %s", certFile)
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Fatalf("Приватный ключ не найден: %s", keyFile)
	}

	// Запуск сервера
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes(), certFile, keyFile); err != nil {
		logrus.Fatalf("Ошибка при запуске HTTPS сервера: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
