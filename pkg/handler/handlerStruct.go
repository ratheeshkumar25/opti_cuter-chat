package handler

import (
	pb "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/proto"
	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/service/interfaces"
)

type ChatServiceServer struct {
	pb.UnimplementedChatServiceServer
	svc interfaces.ChatServiceInter
}

func NewChatServiceServer(svc interfaces.ChatServiceInter) *ChatServiceServer {
	return &ChatServiceServer{
		svc: svc,
	}
}
