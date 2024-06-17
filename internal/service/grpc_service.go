package service

import (
	"context"
	"fmt"
	"log"
	"net"

	"firebase.google.com/go/messaging"
	"github.com/tanush-128/openzo_backend/notification/config"
	"github.com/tanush-128/openzo_backend/notification/internal/pb"
	"github.com/tanush-128/openzo_backend/notification/internal/utils"

	"google.golang.org/grpc"
)

type Server struct {
	pb.NotificationServiceServer
}

func GrpcServer(
	cfg *config.Config,
	server *Server,
) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server listening at %v", lis.Addr())
	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, server)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func (s *Server) SendNotification(ctx context.Context, req *pb.Notification) (*pb.Status, error) {
	// Implement your business logic here
	err := utils.SendNotification(&messaging.Message{
		Notification: &messaging.Notification{
			Title:    req.Title,
			Body:     req.Body,
			ImageURL: req.ImageURL,
		},
		Token: req.Token,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Status{
		Status: "Success",
	}, nil
}

func (s *Server) SendData(ctx context.Context, req *pb.Data) (*pb.Status, error) {
	// Implement your business logic here
	err := utils.SendNotification(&messaging.Message{
		Data:  req.Data,
		Token: req.Token,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Status{
		Status: "Success",
	}, nil
}
