package handler

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/model"
	pb "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/proto"
)

// MessageQueue is a thread-safe channel-based queue for storing messages.
type MessageQueue struct {
	queue chan model.History
}

var messageQueue = MessageQueue{
	queue: make(chan model.History),
}

// Connect establishes a bidirectional streaming connection for chat.
func (c *ChatServiceServer) Connect(csi pb.ChatService_ConnectServer) error {
	errCh := make(chan error)

	// Goroutine to receive messages from the client
	go c.receiveFromStream(csi, errCh)

	// Goroutine to send messages to the client
	go c.sendToStream(csi, errCh)

	// Keep the connection alive or handle the first error
	return <-errCh
}

// receiveFromStream handles incoming messages from the client.
func (c *ChatServiceServer) receiveFromStream(csi pb.ChatService_ConnectServer, errCh chan error) {
	for {
		msg, err := csi.Recv()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("Client disconnected gracefully.")
			} else {
				log.Printf("Error receiving message from client: %v", err)
			}
			errCh <- err
			return
		}

		// Add the received message to the queue
		messageQueue.queue <- model.History{
			UserID:     uint(msg.User_ID),
			ReceiverID: uint(msg.Receiver_ID),
			Message:    msg.Content,
		}

		// Persist the message in the database (business logic)
		go func() {
			if err := c.svc.CreateChatService(msg); err != nil {
				log.Printf("Error saving message to database: %v", err)
			}
		}()
	}
}

// sendToStream sends queued messages to the client.
func (c *ChatServiceServer) sendToStream(csi pb.ChatService_ConnectServer, errCh chan error) {
	select {
	case msg := <-messageQueue.queue:
		if msg.UserID != msg.ReceiverID { // Avoid sending messages to the same user
			// Check if the context is canceled before sending
			if err := csi.Context().Err(); err != nil {
				log.Println("Context canceled, stopping message send.")
				errCh <- err
				return
			}
			err := csi.Send(&pb.Message{
				User_ID:     uint32(msg.UserID),
				Receiver_ID: uint32(msg.ReceiverID),
				Content:     msg.Message,
			})
			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Println("Client disconnected during send.")
				} else {
					log.Printf("Error sending message to client: %v", err)
				}
				errCh <- err
				return
			}
		}
	default:
		// Sleep briefly if no messages to process
		time.Sleep(100 * time.Millisecond)
	}
}

// FetchHistory fetches chat history for a specific conversation.
func (c *ChatServiceServer) FetchHistory(ctx context.Context, p *pb.ChatID) (*pb.ChatHistory, error) {
	response, err := c.svc.FetchChatService(p)
	if err != nil {
		log.Printf("Error fetching chat history: %v", err)
		return nil, err
	}
	return response, nil
}

func (c *ChatServiceServer) StartVideoCall(ctx context.Context, p *pb.VideoCallRequest) (*pb.VideoCallResponse, error) {
	response, err := c.svc.StartVideoCallService(p)
	if err != nil {
		log.Printf("Error in videocall request:%v", err)
		return nil, err
	}
	return response, nil
}

func (c *ChatServiceServer) FetchVideoCall(ctx context.Context, p *pb.ChatID) (*pb.ChatHistory, error) {
	response, err := c.svc.FetchChatService(p)
	if err != nil {
		log.Printf("Error fetching videocall history: %v", err)
		return nil, err
	}
	return response, nil
}
