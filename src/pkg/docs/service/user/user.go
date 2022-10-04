package user

import (
	"fmt"
	"hms/gateway/pkg/docs/service"
	"log"
	"strings"

	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto"
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

func (s *Service) Register(userID, systemID, password, role string) (err error) {
	_, userPrivateKey, err := s.Doc.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("register user error: %w userID %s", err, userID)
	}

	privateUserKey := userPrivateKey[:]
	privateKey, err := crypto.ToECDSA(privateUserKey)
	if err != nil {
		return fmt.Errorf("register user error: %w userID %s", err, userID)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	pwdHash, err := s.generateHash(userID, systemID, password)
	if err != nil {
		return fmt.Errorf("register user error: %w userID %s, password: %s", err, userID, password)
	}

	log.Printf("s.Doc.Infra.Index.userAdd(%s, %s, %s, %s)", address, userID, role, pwdHash)
	//TODO s.Doc.Infra.Index.userAdd(address userAddr, bytes32 id, Role role, bytes pwdHash)
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
