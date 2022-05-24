package services

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"strconv"
)

type CommitThreadService struct {
	CommitThreadRepository interfaces.CommitThreadRepository
}

// TODO: devide in two parts: create thread and create specific thread (like review/request/commit)
func (serv CommitThreadService) StartThread(c *gin.Context, thread models.CommitThreadNoId, aid string, cid string) (models.CommitThread, error) {
	// TODO: check if user is authenticated

	// check model has same aid and cid as params
	intCid, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return models.CommitThread{}, err
	}
	intAid, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		return models.CommitThread{}, err
	}
	if thread.CommitId != intCid || thread.ArticleId != intAid {
		return models.CommitThread{}, errors.New("parameters in url not equal to the thread object")
	}

	// create thread
	threadEntity, err := serv.CommitThreadRepository.CreateThread(aid)
	if err != nil {
		return models.CommitThread{}, err
	}

	// create committhread entity
	threadNoId := models.CommitThreadNoId{
		ArticleId: intAid,
		CommitId:  thread.CommitId,
		ThreadId:  threadEntity.Id,
		Comment:   thread.Comment,
	}
	commitThreadEntity, err := serv.CommitThreadRepository.CreateCommitThread(threadNoId)

	// TODO create and save comment

	// return committhread model or error
	return commitThreadEntity, err

}
