package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/takumi/tailcall-multi-grpc-test/gen/go/name"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

type nameServer struct {
	pb.UnimplementedNameServiceServer
	users map[int32]string
}

func (s *nameServer) GetName(ctx context.Context, req *pb.GetNameRequest) (*pb.GetNameResponse, error) {
	name, ok := s.users[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.Id)
	}
	return &pb.GetNameResponse{Name: name}, nil
}

func loadUsers(filename string) (map[int32]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	userMap := make(map[int32]string)
	for _, user := range users {
		userMap[user.ID] = user.Name
	}

	return userMap, nil
}

func main() {
	users, err := loadUsers("users.json")
	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNameServiceServer(s, &nameServer{users: users})

	log.Printf("Name service listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
