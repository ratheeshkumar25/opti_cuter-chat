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
		select {
		case <-csi.Context().Done():
			log.Println("Client disconnected or context canceled.")
			errCh <- csi.Context().Err()
			return
		default:
			// Proceed to receive the message from the stream.
			msg, err := csi.Recv()
			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Println("Client disconnected gracefully.")
					return
				}
				log.Printf("Error receiving message from client: %v", err)
				errCh <- err
				return
			}

			// Process the message and queue it for later
			messageQueue.queue <- model.History{
				UserID:     uint(msg.UserId),
				ReceiverID: uint(msg.ReceiverId),
				Message:    msg.Content,
			}

			// Persist the message asynchronously in the database
			go func() {
				if err := c.svc.CreateChatService(msg); err != nil {
					log.Printf("Error saving message to database: %v", err)
				}
			}()
		}
	}
}

// sendToStream sends queued messages to the client.
func (c *ChatServiceServer) sendToStream(csi pb.ChatService_ConnectServer, errCh chan error) {
	for {
		select {
		case <-csi.Context().Done():
			log.Println("Client disconnected or context canceled.")
			errCh <- csi.Context().Err()
			return
		default:
			// Try to send a message to the client from the queue.
			select {
			case msg := <-messageQueue.queue:
				//  if the message is for the same user
				if msg.UserID != msg.ReceiverID {
					// Check if the context is canceled before sending the message
					if err := csi.Context().Err(); err != nil {
						log.Println("Context canceled, stopping message send.")
						errCh <- err
						return
					}

					// Send the message to the client
					err := csi.Send(&pb.Message{
						UserId:     uint32(msg.UserID),
						ReceiverId: uint32(msg.ReceiverID),
						Content:    msg.Message,
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
				} else {
					// Send self-message back to the sender
					err := csi.Send(&pb.Message{
						UserId:     uint32(msg.UserID),
						ReceiverId: uint32(msg.UserID),
						Content:    msg.Message,
					})
					if err != nil {
						log.Printf("Error sending self-message to client: %v", err)
						errCh <- err
					}
				}
			default:
				// Sleep briefly if no messages to process
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (c *ChatServiceServer) FetchHistory(ctx context.Context, p *pb.ChatID) (*pb.ChatHistory, error) {
	response, err := c.svc.FetchChatService(p)
	if err != nil {
		return response, err
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
