package handler

import (
	"context"
	"log"

	pb "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/proto"
)

func (c *ChatServiceServer) SubmitReview(ctx context.Context, p *pb.ReviewRequest) (*pb.ReviewResponse, error) {
	response, err := c.svc.SubmitReviewService(p)
	if err != nil {
		log.Printf("Error in submiting review request:%v", err)
		return nil, err
	}
	return response, nil
}

func (c *ChatServiceServer) FetchReviews(ctx context.Context, p *pb.ChatMaterialID) (*pb.ReviewList, error) {
	response, err := c.svc.FetchReviewService(p)
	if err != nil {
		log.Printf("Error in fetching review request:%v", err)
		return nil, err
	}
	return response, nil
}

func (c *ChatServiceServer) AddVideoChunk(ctx context.Context, p *pb.VideoUploadRequest) (*pb.VideoUploadResponse, error) {
	response, err := c.svc.UploadVideoService(p)
	if err != nil {
		log.Printf("Error in uploading video request:%v", err)
		return nil, err
	}
	return response, nil
}

func (c *ChatServiceServer) FetchVideos(ctx context.Context, p *pb.FetchVideoRequest) (*pb.FetchVideoResponse, error) {
	response, err := c.svc.FetchVideoService(p)
	if err != nil {
		log.Printf("Error in uploading video request:%v", err)
		return nil, err
	}
	return response, nil
}
