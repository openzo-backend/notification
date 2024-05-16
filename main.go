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

	// // Initialize gRPC server
	// grpcServer := grpc.NewServer()
	// notificationpb.RegisternotificationServiceServer(grpcServer, service.NewGrpcnotificationService(notificationRepository, notificationService))
	// reflection.Register(grpcServer) // Optional for server reflection

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
