package server

import (
	"fmt"
	"log"
	"net"

	pb "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/proto"
	"google.golang.org/grpc"
)

func NewGrpcUserServer(port string, handler pb.ChatServiceServer) error {
	log.Println("connecting to gRPC server")
	addr := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("error creating listener on %v", addr)
		return err
	}
	grpc := grpc.NewServer()
	pb.RegisterChatServiceServer(grpc, handler)

	log.Printf("listening on gRPC server %v", port)
	err = grpc.Serve(lis)
	if err != nil {
		log.Println("error connecting to gRPC server")
		return err
	}
	return nil
}
