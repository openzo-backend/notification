package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/notification/config"
	handlers "github.com/tanush-128/openzo_backend/notification/internal/api"
	"github.com/tanush-128/openzo_backend/notification/internal/middlewares"
	"github.com/tanush-128/openzo_backend/notification/internal/pb"
	"github.com/tanush-128/openzo_backend/notification/internal/repository"
	"github.com/tanush-128/openzo_backend/notification/internal/service"
	"google.golang.org/grpc"
)

var UserClient pb.UserServiceClient

type Notification struct {
	Message  string `json:"message"`
	FCMToken string `json:"fcm_token"`
	Data     string `json:"data,omitempty"`
	Topic    string `json:"topic,omitempty"`
}

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}

	// Connect to database
	db, err := connectToDB(cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}

	// Set up notification repository and service
	notificationRepository := repository.NewLocalNotificationRepository(db)
	notificationService := service.NewLocalNotificationService(notificationRepository)

	// Start Kafka consumer in a goroutine
	go consumeKafka()

	// Initialize gRPC client
	conn, err := grpc.Dial(cfg.UserGrpc, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to gRPC: %v", err)
	}
	defer conn.Close()
	UserClient = pb.NewUserServiceClient(conn)

	// Start gRPC server
	go service.GrpcServer(cfg, &service.Server{})

	// Initialize HTTP server with Gin
	router := gin.Default()
	handler := handlers.NewHandler(&notificationService)

	// Define routes
	router.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/", handler.CreateNotification)
	router.POST("pincode/", handler.CreateLocalNotification)
	router.GET("store/:id", handler.GetNotificationsByStoreID)
	router.GET("/:id", handler.GetNotificationByID)
	router.GET("pincode/:pincode", handler.GetNotifications)
	router.Use(middlewares.NewMiddleware(UserClient).JwtMiddleware)
	router.DELETE("/:id", handler.DeleteNotification)

	// Start HTTP server
	router.Run(fmt.Sprintf(":%s", cfg.HTTPPort))
}

// consumeKafka listens to Kafka notifications and processes them

// ReadConfig reads the Kafka configuration from a properties file
