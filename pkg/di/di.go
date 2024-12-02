package di

import (
	"log"

	"github.com/ratheeshkumar25/opti_cut_chat_service/config"
	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/db"
	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/handler"
	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/repo"
	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/server"
	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/service"
)

func Init() {
	cnfg := config.LoadConfig()

	db, err := db.ConnectMongoDB(cnfg)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	mongoDB := db.Database(cnfg.DBName)
	chatRepo := repo.NewChatRepository(mongoDB)

	chatSVC := service.NewChatService(chatRepo)
	chatServer := handler.NewChatServiceServer(chatSVC)
	err = server.NewGrpcUserServer(cnfg.GrpcPort, chatServer)
	if err != nil {
		log.Fatalf("failed to start gRPC server %v", err.Error())
	}
}
