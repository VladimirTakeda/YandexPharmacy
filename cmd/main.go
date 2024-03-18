package main

import (
	"context"
	"flag"
	"fmt"
	"git.yandex-academy.ru/ooornament/code_architecture/cmd/internal"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/config"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/db/postgresql"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/handler"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/repository"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/service"
	"github.com/antelman107/net-wait-go/wait"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "deploy", "config folder")
	flag.Parse()
}

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config %s", err.Error())
		os.Exit(1)
	}

	ctx := context.Background()
	var cfg config.Config
	err := confita.NewLoader(
		file.NewBackend(fmt.Sprintf("%s/default.yaml", configPath)),
	).Load(ctx, &cfg)
	if err != nil {
		fmt.Printf("failed to parse config: %s\n", err.Error())
		return
	}

	if !wait.New(
		wait.WithProto("tcp"),
		wait.WithWait(200*time.Millisecond),
		wait.WithBreak(50*time.Millisecond),
		wait.WithDeadline(15*time.Second),
		wait.WithDebug(true),
	).Do([]string{"localhost:5432"}) {
		logrus.Fatalf("db is not available")
		return
	}

	postgres, err := postgresql.NewPostgresDb(cfg.Postgres)
	if err != nil {
		fmt.Printf("failed to connect postgresql: %s\n", err.Error())
		return
	}

	repos := repository.NewRepository(postgres)
	services := service.NewCartCheckService(repos)
	handlers := handler.NewHandler(services)

	srv := new(internal.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.SetupRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
			os.Exit(1)
		}
	}()

	logrus.Print("App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("cmd/internal")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
