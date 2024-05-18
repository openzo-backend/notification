package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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
}

func consumeKafka() {
	conf := ReadConfig()

	topic := "notification"

	// sets the consumer group ID and offset
	conf["group.id"] = "go-group-1"
	conf["auto.offset.reset"] = "earliest"

	// creates a new consumer and subscribes to your topic
	consumer, _ := kafka.NewConsumer(&conf)
	consumer.SubscribeTopics([]string{topic}, nil)
	var notification Notification
	run := true
	for run {
		// consumes messages from the subscribed topic and prints them to the console
		e := consumer.Poll(1000)
		switch ev := e.(type) {
		case *kafka.Message:
			// application-specific processing
			err := json.Unmarshal(ev.Value, &notification)
			if err != nil {
				fmt.Println("Error unmarshalling JSON: ", err)
			}
			fmt.Println("Sending notification", notification)

			err = utils.SendNotification(&messaging.Message{
				Notification: &messaging.Notification{
					Title: "New Order!",
					Body:  notification.Message,
				},
				Token: notification.FCMToken,
			})

			if err != nil {
				fmt.Println("Error sending notification: ", err)
			}

		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", ev)
			run = false
		}
	}

	// closes the consumer connection
	consumer.Close()

}
