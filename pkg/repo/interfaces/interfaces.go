package interfaces

import "github.com/ratheeshkumar25/opti_cut_chat_service/pkg/model"

type ChatRepoInter interface {
	Createchat(chat *model.History) error
	Findchat(userID, receiverID uint) (*[]model.History, error)
	// CreateVideoCallLog(call *model.VideoCallLog) error
	LogVideoCall(videoCall *model.VideoCall) error
	GetVideoCallHistory(userID, receiverID uint) ([]model.VideoCall, error)
}
