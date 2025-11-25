package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/takumi/tailcall-multi-grpc-test/gen/go/age"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

type ageServer struct {
	pb.UnimplementedAgeServiceServer
	users map[int32]int32
}

func (s *ageServer) GetAge(ctx context.Context, req *pb.GetAgeRequest) (*pb.GetAgeResponse, error) {
	age, ok := s.users[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.Id)
	}
	return &pb.GetAgeResponse{Age: age}, nil
}

func loadUsers(filename string) (map[int32]int32, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	userMap := make(map[int32]int32)
	for _, user := range users {
		userMap[user.ID] = user.Age
	}

	return userMap, nil
}

func main() {
	users, err := loadUsers("users.json")
	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAgeServiceServer(s, &ageServer{users: users})

	log.Printf("Age service listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
