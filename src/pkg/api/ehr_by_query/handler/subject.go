package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/types"
	"log"
)

type Subject struct{}

func (s *Subject) CanProcess(c *gin.Context) bool {
	if "" != c.Query("subject_id") && "" != c.Query("namespace") {
		return true
	}
	return false
}

// Handle EHR searching by subject_id and subject namespace
func (s *Subject) Handle(c *gin.Context) ([]byte, error) {
	subjectId := c.Query("subject_id")
	namespace := c.Query("namespace")

	docService := service.NewDefaultDocumentService()

	userId := c.GetString("userId")
	if userId == "" {
		return nil, errors.New("userId is empty")
	}

	ehrId, err := docService.SubjectIndex.GetEhrBySubject(subjectId, namespace)
	if err != nil {
		return nil, err
	}

	// Getting docStorageId
	doc, err := docService.DocsIndex.GetLastByType(ehrId, types.EHR)
	if err != nil {
		log.Println("GetLastDocIndexByType", "ehrId", ehrId, err)
		return nil, err
	}

	// Getting doc from storage
	docDecrypted, err := docService.GetDocFromStorageById(userId, doc.StorageId, []byte(ehrId))
	if err != nil {
		return nil, err
	}

	return docDecrypted, nil
}
