package service

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/model"
	pb "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/proto"
	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/repo/interfaces"
	inter "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/service/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type chatService struct {
	repo interfaces.ChatRepoInter
}

func generateJitsiRoomURL(userID, receiverID uint) string {
	// Combine userID and receiverID to form a unique room name
	roomName := fmt.Sprintf("call-%d-%d", userID, receiverID)

	// Create the base Jitsi Meet URL
	baseURL := "https://meet.jit.si/" + roomName

	// Add query parameters for user and receiver IDs
	params := url.Values{}
	params.Add("user_id", fmt.Sprintf("%d", userID))
	params.Add("receiver_id", fmt.Sprintf("%d", receiverID))

	// Return the full Jitsi room URL with query parameters
	return baseURL + "?" + params.Encode()
}

// FetchVideoCallService implements interfaces.ChatServiceInter.
func (c *chatService) FetchVideoCallService(p *pb.ChatID) (*pb.ChatHistory, error) {
	//retrieve video call history
	calls, err := c.repo.GetVideoCallHistory(uint(p.User_ID), uint(p.Receiver_ID))
	if err != nil {
		return nil, err
	}
	var history pb.ChatHistory
	for _, call := range calls {
		history.Chats = append(history.Chats, &pb.Message{
			User_ID:     uint32(call.UserID),
			Receiver_ID: uint32(call.ReceiverID),
			Content:     call.RoomURL,
		})
	}
	return &history, nil
}

// StartVideoCall generates a new video call URL using Jitsi and stores it
func (c *chatService) StartVideoCallService(p *pb.VideoCallRequest) (*pb.VideoCallResponse, error) {
	roomURL := generateJitsiRoomURL(uint(p.User_ID), uint(p.Receiver_ID))
	videoCall := &model.VideoCall{
		UserID:     uint(p.User_ID),
		ReceiverID: uint(p.Receiver_ID),
		RoomURL:    roomURL,
		Timestamp:  primitive.NewDateTimeFromTime(time.Now()),
	}
	//store video call information to db
	if err := c.repo.LogVideoCall(videoCall); err != nil {
		return nil, err
	}

	return &pb.VideoCallResponse{RoomUrl: roomURL}, nil
}

func NewChatService(repo interfaces.ChatRepoInter) inter.ChatServiceInter {
	return &chatService{
		repo: repo,
	}
}