package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	orders "github.com/samarec1812/wb-orders-test0"
	"github.com/samarec1812/wb-orders-test0/pkg/handler"
	"github.com/samarec1812/wb-orders-test0/pkg/repository"
	"github.com/samarec1812/wb-orders-test0/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	Repos *repository.Repository
)

func stanSubscribe() {

	sc, _ := stan.Connect("prod", "sub")
	defer sc.Close()

	_, err := sc.Subscribe("samarec", func(m *stan.Msg) {
		var order orders.Orders
		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			logrus.Fatalf("error model")
		} else {
			url := "http://127.0.0.1:8089/api/send"
			_, err = http.Post(url, "application/json", bytes.NewBuffer(m.Data))
			if err != nil {
				logrus.Fatalf("error create %s", err.Error())
			}
		}
	}, stan.AckWait(5* time.Second), stan.DurableName("sub-1"))

	if err != nil {
		logrus.Fatalf("error with subscribe %s", err.Error())
		return
	}

	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
	//runtime.Goexit()
}


func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName: viper.GetString("db.dbname"),
		SSLMode: viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	rd, err := repository.NewRedisDB(repository.ConfigRedis{
		Addr: viper.GetString("rd.addr"),
		Password: os.Getenv("RD_PASSWORD"),
		DB: viper.Get("rd.db").(int),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize redis: %s", err.Error())
	}

	repos := repository.NewRepository(db, rd)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)



	srv := new(orders.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	// wait signal to shutdown server with a timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Println("Shutting down server. ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Errorf("error occured on db connection close: %s", err.Error())
		}
		if err := rd.Close(); err != nil {
			logrus.Errorf("error ocurred on rd connection close: %s", err.Error())
		}
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %s", err.Error())
	}

	logrus.Println("Server exiting")
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}