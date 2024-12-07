package service

import (
	"sort"

	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/model"
	pb "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/proto"
)

func (c *chatService) CreateChatService(p *pb.Message) error {
	chat := &model.History{
		UserID:     uint(p.UserId),
		ReceiverID: uint(p.ReceiverId),
		Message:    p.Content,
	}
	err := c.repo.Createchat(chat)
	if err != nil {
		return err
	}

	return nil
}

// FetchChatService implements interfaces.ChatServiceInter.
func (c *chatService) FetchChatService(p *pb.ChatID) (*pb.ChatHistory, error) {

	userHistory, err := c.repo.Findchat(uint(p.UserId), uint(p.ReceiverId))
	if err != nil {
		return nil, err
	}
	receiverHistory, err := c.repo.Findchat(uint(p.ReceiverId), uint(p.UserId))
	if err != nil {
		return nil, err
	}
	var chats []*pb.Message
	for _, msg := range *userHistory {
		chats = append(chats, &pb.Message{
			ChatId:     uint32(msg.ID.Timestamp().Unix()),
			UserId:     uint32(msg.UserID),
			ReceiverId: uint32(msg.ReceiverID),
			Content:    msg.Message,
		})

	}
	for _, msg := range *receiverHistory {
		chats = append(chats, &pb.Message{
			ChatId:     uint32(msg.ID.Timestamp().Unix()),
			UserId:     uint32(msg.UserID),
			ReceiverId: uint32(msg.ReceiverID),
			Content:    msg.Message,
		})
	}
	sortByChatID(chats)
	return &pb.ChatHistory{
		Chats: chats,
	}, nil
}

func sortByChatID(chats []*pb.Message) {
	sort.Slice(chats, func(i, j int) bool {
		return chats[i].ChatId < chats[j].ChatId
	})
}
