package service

import (
	"context"
	"fmt"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	proc "hms/gateway/pkg/docs/service/processing"
	"strings"

	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
)

type Service struct {
	Doc *service.DefaultDocumentService
}

const (
	N            = 1048576
	r            = 8
	p            = 1
	keyLen       = 32
	saltLenBytes = 16
)

func NewUserService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		docService,
	}
}

func (s *Service) Register(ctx context.Context, procRequest *proc.Request, user *model.UserCreateRequest) (err error) {
	_, userPrivateKey, err := s.Doc.Infra.Keystore.Get(user.UserID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, user.UserID)
	}

	privateUserKey := userPrivateKey[:]
	privateKey, err := crypto.ToECDSA(privateUserKey)

	if err != nil {
		return fmt.Errorf("crypto.ToECDSA error: %w userID %s", err, user.UserID)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	ehrSystemID := ctx.(*gin.Context).GetString("ehrSystemID")
	pwdHash, err := s.generateHash(user.UserID, ehrSystemID, user.Password)

	if err != nil {
		return fmt.Errorf("register s.generateHash error: %w userID %s, password: %s", err, user.UserID, user.Password)
	}

	requestID := ctx.(*gin.Context).GetString("reqId")
	txHash, err := s.Doc.Infra.Index.UserAdd(requestID, address, user.UserID, user.Role, pwdHash)

	if err != nil {
		return fmt.Errorf("Index.DeleteDoc error: %w", err)
	}

	procRequest.AddEthereumTx(proc.TxDeleteDoc, txHash)

	return nil
}

func (s *Service) generateHash(phrases ...string) ([]byte, error) {
	hash := strings.Join(phrases, "")

	salt := make([]byte, saltLenBytes)

	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	pwdHash, err := scrypt.Key([]byte(hash), salt, N, r, p, keyLen)
	if err != nil {
		return nil, err
	}

	return append(salt, pwdHash...), nil
}
