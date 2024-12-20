package interfaces

import (
	pb "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/proto"
)

type ChatServiceInter interface {
	CreateChatService(p *pb.Message) error
	FetchChatService(p *pb.ChatID) (*pb.ChatHistory, error)
	StartVideoCallService(p *pb.VideoCallRequest) (*pb.VideoCallResponse, error)
	FetchVideoCallService(p *pb.ChatID) (*pb.ChatHistory, error)
	SubmitReviewService(p *pb.ReviewRequest) (*pb.ReviewResponse, error)
	FetchReviewService(p *pb.ChatMaterialID) (*pb.ReviewList, error)
	UploadVideoService(p *pb.VideoUploadRequest) (*pb.VideoUploadResponse, error)
	FetchVideoService(p *pb.FetchVideoRequest) (*pb.FetchVideoResponse, error)
}
