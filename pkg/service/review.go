package service

import (
	"context"
	"fmt"
	"hash/crc32"
	"time"

	materialpb "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/client/material/pb"
	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/model"
	pb "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubmitReviewService implements interfaces.ChatServiceInter.
func (c *chatService) SubmitReviewService(p *pb.ReviewRequest) (*pb.ReviewResponse, error) {
	materialIDResp, err := c.MaterialClient.FindMaterialByID(context.Background(), &materialpb.MaterialID{ID: p.MaterialId})
	if err != nil || materialIDResp == nil {
		return &pb.ReviewResponse{
			Status:  pb.ReviewResponse_FAILED,
			Message: "materila not found or service unavailable",
		}, fmt.Errorf("failed to validate materialID %v", err)
	}

	review := &model.Review{
		ID:         primitive.NewObjectID(),
		UserID:     p.UserId,
		MaterialID: p.MaterialId,
		Content:    p.ReviewText,
		Rating:     uint32(p.Rating),
		Timestamp:  time.Now(),
	}

	if err := c.repo.AddReview(review); err != nil {
		return &pb.ReviewResponse{
			Status:  pb.ReviewResponse_FAILED,
			Message: "failed to save review",
		}, fmt.Errorf("error in saving review %v", err)
	}
	return &pb.ReviewResponse{
		Status:  pb.ReviewResponse_SUCCESS,
		Message: "review submitted successfully",
	}, nil
}

// FetchReviewService implements interfaces.ChatServiceInter.
func (c *chatService) FetchReviewService(p *pb.ChatMaterialID) (*pb.ReviewList, error) {
	// Fetch reviews from the repository
	reviews, err := c.repo.FindReviewMaterialID(uint(p.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	// Convert reviews to protobuf format
	var reviewList []*pb.Review
	for _, review := range reviews {
		reviewList = append(reviewList, &pb.Review{
			ReviewId:   hashObjectIDToUint32(review.ID),
			UserId:     review.UserID,
			MaterialId: review.MaterialID,
			ReviewText: review.Content,
			Rating:     int32(review.Rating),
			Timestamp:  review.Timestamp.Format(time.RFC3339),
		})
	}

	return &pb.ReviewList{Reviews: reviewList}, nil
}

// UploadVideoService handles video chunk uploads.
func (c *chatService) UploadVideoService(p *pb.VideoUploadRequest) (*pb.VideoUploadResponse, error) {
	var videoID primitive.ObjectID

	if p.IsFirstChunk {
		// Generate a new VideoID for the first chunk
		videoID = primitive.NewObjectID()

		// Log the generated VideoID
		fmt.Printf("Generated new VideoId: %s\n", videoID.Hex())

		// Create a new Video document with metadata
		video := &model.Video{
			ID:         videoID,
			UserID:     p.UserId,
			MaterialID: p.MaterialId,
			FileName:   p.FileName,
			VideoURL:   p.VideoUrl,
			Timestamp:  time.Now(),
			Chunks:     []model.VideoChunk{}, // Initialize empty chunks array
		}

		// Save the video metadata in the database
		if err := c.repo.AddVideo(video); err != nil {
			return &pb.VideoUploadResponse{
				Status:  pb.VideoUploadResponse_FAILED,
				Message: "Failed to create video metadata",
			}, fmt.Errorf("error creating video metadata: %v", err)
		}
	} else {
		// Log the received VideoId for debugging
		fmt.Printf("Received VideoId: %s\n", p.VideoId)

		// Validate VideoId length
		if len(p.VideoId) != 24 {
			return &pb.VideoUploadResponse{
				Status:  pb.VideoUploadResponse_FAILED,
				Message: "Invalid video ID: length must be 24 characters",
			}, fmt.Errorf("invalid video ID: length must be 24 characters, got: %d", len(p.VideoId))
		}

		// Parse VideoID for subsequent chunks
		var err error
		videoID, err = primitive.ObjectIDFromHex(p.VideoId)
		if err != nil {
			return &pb.VideoUploadResponse{
				Status:  pb.VideoUploadResponse_FAILED,
				Message: fmt.Sprintf("Invalid video ID: %s", p.VideoId),
			}, fmt.Errorf("invalid video ID: %s, error: %v", p.VideoId, err)
		}
	}

	// Add the video chunk to the database
	chunk := model.VideoChunk{
		ChunkID:    primitive.NewObjectID(),
		ChunkData:  p.VideoData,
		ChunkOrder: int(p.ChunkOrder),
	}

	if err := c.repo.AddVideoChunk(videoID, &chunk); err != nil {
		return &pb.VideoUploadResponse{
			Status:  pb.VideoUploadResponse_FAILED,
			Message: "Failed to save video chunk",
		}, fmt.Errorf("error saving video chunk: %v", err)
	}

	// Finalize the video if it's the last chunk
	if p.IsLastChunk {
		if err := c.repo.FinalizeVideo(videoID, p.FileName, p.UserId); err != nil {
			return &pb.VideoUploadResponse{
				Status:  pb.VideoUploadResponse_FAILED,
				Message: "Failed to finalize video upload",
			}, fmt.Errorf("error finalizing video: %v", err)
		}
	}

	// Respond with the VideoID
	return &pb.VideoUploadResponse{
		Status:  pb.VideoUploadResponse_SUCCESS,
		Message: "Video chunk uploaded successfully",
		VideoId: videoID.Hex(),
	}, nil
}

// FetchVideoService implements interfaces.ChatServiceInter.
func (c *chatService) FetchVideoService(p *pb.FetchVideoRequest) (*pb.FetchVideoResponse, error) {
	videos, err := c.repo.FindVideoByMaterialID(uint(p.MaterialId))
	if err != nil {
		return &pb.FetchVideoResponse{
			Videos: nil,
		}, fmt.Errorf("error fetching videos: %v", err)
	}

	// Map the video objects into the protobuf response format
	var videoMetadataList []*pb.VideoMetadata
	for _, video := range videos {
		videoMetadata := &pb.VideoMetadata{
			VideoId:    video.ID.Hex(),
			MaterialId: video.MaterialID,
			UserId:     video.UserID,
			FileName:   video.FileName,
			VideoUrl:   video.VideoURL,
			Timestamp:  video.Timestamp.Format(time.RFC3339),
		}
		videoMetadataList = append(videoMetadataList, videoMetadata)
	}

	// Return the list of video metadata in the response
	return &pb.FetchVideoResponse{
		Videos: videoMetadataList,
	}, nil
}

// Hash the primitiveobjectID
func hashObjectIDToUint32(objectID primitive.ObjectID) uint32 {
	hash := crc32.ChecksumIEEE([]byte(objectID.Hex()))
	return hash
}

// func (c *chatService) UploadVideoService(p *pb.VideoUploadRequest) (*pb.VideoUploadResponse, error) {
// 	var videoID primitive.ObjectID
// 	if p.IsFirstChunk {
// 		// Generate a new video ID for the first chunk
// 		videoID = primitive.NewObjectID()

// 		// Initialize video metadata
// 		video := &model.Video{
// 			ID:         videoID,
// 			UserID:     p.UserId,
// 			MaterialID: p.MaterialId,
// 			FileName:   p.FileName,
// 			VideoURL:   p.VideoUrl,
// 			Timestamp:  time.Now(),
// 			Chunks:     []model.VideoChunk{}, // Empty initially
// 		}

// 		// Add metadata to the repository
// 		if err := c.repo.AddVideo(video); err != nil {
// 			return &pb.VideoUploadResponse{
// 				Status:  pb.VideoUploadResponse_FAILED,
// 				Message: "Failed to add video metadata",
// 			}, fmt.Errorf("error adding video metadata: %v", err)
// 		}
// 	} else {
// 		// Parse video ID for subsequent chunks
// 		var err error
// 		videoID, err = primitive.ObjectIDFromHex(p.VideoId)
// 		if err != nil {
// 			return &pb.VideoUploadResponse{
// 				Status:  pb.VideoUploadResponse_FAILED,
// 				Message: "Invalid video ID",
// 			}, fmt.Errorf("invalid video ID: %v", err)
// 		}
// 	}

// 	// Create a chunk for this part of the video
// 	chunk := &model.VideoChunk{
// 		ChunkID:    primitive.NewObjectID(),
// 		ChunkData:  p.VideoData,
// 		ChunkOrder: int(p.ChunkOrder),
// 	}

// 	// Add the chunk to the repository
// 	if err := c.repo.AddVideoChunk(videoID, chunk); err != nil {
// 		return &pb.VideoUploadResponse{
// 			Status:  pb.VideoUploadResponse_FAILED,
// 			Message: "Failed to upload video chunk",
// 		}, fmt.Errorf("error saving video chunk: %v", err)
// 	}

// 	// Finalize the video metadata if it's the last chunk
// 	if p.IsLastChunk {
// 		if err := c.repo.FinalizeVideo(videoID, p.FileName, p.UserId); err != nil {
// 			return &pb.VideoUploadResponse{
// 				Status:  pb.VideoUploadResponse_FAILED,
// 				Message: "Failed to finalize video upload",
// 			}, fmt.Errorf("error finalizing video: %v", err)
// 		}
// 	}

// 	return &pb.VideoUploadResponse{
// 		Status:  pb.VideoUploadResponse_SUCCESS,
// 		Message: "Video chunk uploaded successfully",
// 		VideoId: videoID.Hex(),
// 	}, nil
// }
