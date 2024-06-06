package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/notification/config"
	handlers "github.com/tanush-128/openzo_backend/notification/internal/api"
	"github.com/tanush-128/openzo_backend/notification/internal/middlewares"
	"github.com/tanush-128/openzo_backend/notification/internal/pb"
	"github.com/tanush-128/openzo_backend/notification/internal/repository"
	"github.com/tanush-128/openzo_backend/notification/internal/service"
	"github.com/tanush-128/openzo_backend/notification/internal/utils"
	"google.golang.org/grpc"
)

var UserClient pb.UserServiceClient

type User2 struct {
}

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}

	db, err := connectToDB(cfg) // Implement database connection logic
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}

	notificationRepository := repository.NewLocalNotificationRepository(db)
	notificationService := service.NewLocalNotificationService(notificationRepository)

	go consumeKafka()

	//Initialize gRPC client
	conn, err := grpc.Dial(cfg.UserGrpc, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)
	UserClient = c

	go service.GrpcServer(cfg, &service.Server{})
	// Initialize HTTP server with Gin
	router := gin.Default()
	handler := handlers.NewHandler(&notificationService)

	router.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// router.Use(middlewares.JwtMiddleware(c))
	router.POST("/notifications", handler.CreateNotification)
	router.GET("/notifications/store/:id", handler.GetNotificationsByStoreID)
	router.GET("/notifications/pincode/:pincode", handler.GetNotifications)
	router.GET("/notifications/:id", handler.GetNotificationByID)
	router.Use(middlewares.NewMiddleware(c).JwtMiddleware)
	router.DELETE("/notifications/:id", handler.DeleteNotification)

	// router.Use(middlewares.JwtMiddleware)

	router.Run(fmt.Sprintf(":%s", cfg.HTTPPort))

}

type Notification struct {
	Message  string `json:"message"`
	FCMToken string `json:"fcm_token"`
	Data     string `json:"data,omitempty"`
	Topic    string `json:"topic,omitempty"`
}

func consumeKafka() {
	conf := ReadConfig()
	topic := "notification"
	conf["group.id"] = "go-group-1"
	conf["auto.offset.reset"] = "earliest"

	var notification Notification

	for {
		consumer, err := kafka.NewConsumer(&conf)
		if err != nil {
			log.Printf("Error creating consumer: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = consumer.SubscribeTopics([]string{topic}, nil)
		if err != nil {
			log.Printf("Error subscribing to topic: %v. Retrying in 5 seconds...", err)
			consumer.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		run := true
		for run {
			e := consumer.Poll(1000)
			switch ev := e.(type) {
			case *kafka.Message:
				err := json.Unmarshal(ev.Value, &notification)
				if err != nil {
					log.Printf("Error unmarshalling JSON: %v", err)
					continue
				}

				notifDataMap := map[string]string{}
				err = json.Unmarshal([]byte(notification.Data), &notifDataMap)
				if err != nil {
					log.Printf("Error unmarshalling JSON: %v", err)
					continue
				}
				log.Printf("Notification received: %+v", notification)
				log.Println("Sending notification", notification)

				err = utils.SendNotification(&messaging.Message{
					Notification: &messaging.Notification{
						Title: notification.Topic + " Notification",
						Body:  notification.Message,
					},
					Data:  notifDataMap,
					Topic: notification.Topic,
					Token: notification.FCMToken,
				})

				if err != nil {
					log.Printf("Error sending notification: %v", err)
				}

			case kafka.Error:
				log.Printf("Error: %v", ev)
				run = false
			}
		}

		log.Println("Consumer disconnected. Reconnecting in 5 seconds...")
		consumer.Close()
		time.Sleep(5 * time.Second)
	}
}
