package main

import (
	pb "chat-app/chat"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type chatServer struct {
	pb.UnimplementedChatServiceServer
	Clients map[string]chan *pb.Message
}

func (s *chatServer) Connect(client *pb.Client, stream pb.ChatService_ConnectServer) error {
	s.Clients[client.Id] = make(chan *pb.Message)
	for {
		message := <-s.Clients[client.Id]
		if err := stream.Send(message); err != nil {
			return err
		}
	}
	return nil
}

func (s *chatServer) SendMessage(ctx context.Context, message *pb.Message) (*emptypb.Empty, error) {
	for _, c := range s.Clients {
		c <- message
	}
	return &emptypb.Empty{}, nil
}

func newServer() *chatServer {
	return &chatServer{
		Clients: make(map[string]chan *pb.Message),
	}
}

func main() {
	listener, error := net.Listen("tcp", ":8080")
	if error != nil {
		log.Fatalf("Failed to listen: %v", error)
	}
	fmt.Println("Server is running on port 8080")

	// Create a new gRPC server and register the chat service
	var options []grpc.ServerOption
	grpcServer := grpc.NewServer(options...)
	pb.RegisterChatServiceServer(grpcServer, newServer())
	grpcServer.Serve(listener)
}
