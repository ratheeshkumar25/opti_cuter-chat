package interfaces

import (
	"github.com/ratheeshkumar25/opti_cut_chat_service/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRepoInter interface {
	Createchat(chat *model.History) error
	Findchat(userID, receiverID uint) (*[]model.History, error)

	LogVideoCall(videoCall *model.VideoCall) error
	GetVideoCallHistory(userID, receiverID uint) ([]model.VideoCall, error)

	AddReview(review *model.Review) error
	FindReviewMaterialID(materialID uint) ([]model.Review, error)
	AddVideo(video *model.Video) error
	FindVideoByMaterialID(materilaID uint) ([]model.Video, error)
	AddVideoChunk(videoID primitive.ObjectID, chunk *model.VideoChunk) error
	FinalizeVideo(videoID primitive.ObjectID, fileName string, userID uint32) error
	FetchVideoChunks(videoID primitive.ObjectID) ([]model.VideoChunk, error)
}
