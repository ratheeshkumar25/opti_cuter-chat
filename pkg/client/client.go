package client

import (
	"fmt"
	"log"

	"github.com/ratheeshkumar25/opti_cut_chat_service/config"
	pb "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/client/material/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ClientDial(cfg config.Config) (pb.MaterialServiceClient, error) {
	grpcAddr := fmt.Sprintf("material-service:%s", cfg.MaterialPort)
	grpc, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Error dialing to grpc to client : %s", err.Error())
		return nil, err
	}
	log.Printf("Successfully connected to material client at port : %s", cfg.MaterialPort)
	return pb.NewMaterialServiceClient(grpc), nil
}
