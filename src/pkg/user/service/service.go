package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/akyoto/cache"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service/processing"
	proc "hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/infrastructure"
)

type Service struct {
	Infra *infrastructure.Infra
	Proc  *processing.Proc
	Cache *cache.Cache
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

type TokenType uint8

type TokenClaims struct {
	TokenType TokenType `json:"type"`
	jwt.StandardClaims
}

const (
	TokenAccessType TokenType = iota
	TokenRefreshType
)

func NewUserService(cfg *config.Config, infra *infrastructure.Infra, p *processing.Proc) *Service {
	return &Service{
		Infra: infra,
		Proc:  p,
		Cache: cache.New(common.CacheCleanerTimeout),
	}
}

func (s *Service) Register(ctx context.Context, procRequest *proc.Request, user *model.UserCreateRequest, systemID string) (err error) {
	_, userPrivateKey, err := s.Infra.Keystore.Get(user.UserID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, user.UserID)
	}

	pwdHash, err := generateHashFromPassword(systemID, user.UserID, user.Password)
	if err != nil {
		return fmt.Errorf("generateHashFromPassword error: %w", err)
	}

	txHash, err := s.Infra.Index.UserNew(ctx, user.UserID, systemID, user.Role, pwdHash, userPrivateKey, nil)
	if err != nil {
		if errors.Is(err, errors.ErrAlreadyExist) {
			return err
		}

		return fmt.Errorf("Index.UserAdd error: %w", err)
	}

	procRequest.AddEthereumTx(proc.TxUserRegister, txHash, true)

	return nil
}

func (s *Service) Login(ctx context.Context, userID, systemID, password string) (err error) {
	address, err := s.getUserAddress(userID)
	if err != nil {
		return fmt.Errorf("Login s.getUserAddress error: %w", err)
	}

	log.Println("GetUserPasswordHash start")

	pwdHash, err := s.Infra.Index.GetUserPasswordHash(ctx, address)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return err
		}
		return fmt.Errorf("Login.GetUserPasswordHash error: %w", err)
	}

	log.Println("GetUserPasswordHash finish")

	match, err := verifyPassphrase(userID+systemID+password, pwdHash)
	if err != nil {
		return fmt.Errorf("verifyPassphrase error: %w", err)
	}

	if !match {
		return errors.ErrAuthorization
	}

	return nil
}

func (s *Service) getUserAddress(userID string) (eth_common.Address, error) {
	_, userPrivateKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return eth_common.Address{}, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	privateKey, err := crypto.ToECDSA(userPrivateKey[:])
	if err != nil {
		return eth_common.Address{}, fmt.Errorf("crypto.ToECDSA error: %w userID %s", err, userID)
	}

	return crypto.PubkeyToAddress(privateKey.PublicKey), nil
}

func generateHashFromPassword(systemID, userID, password string) ([]byte, error) {
	salt := make([]byte, common.ScryptSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("rand.Read error: %w", err)
	}

	password = userID + systemID + password

	pwdHash, err := scrypt.Key([]byte(password), salt, common.ScryptN, common.ScryptR, common.ScryptP, common.ScryptKeyLen)
	if err != nil {
		return nil, fmt.Errorf("generateHash scrypt.Key error: %w", err)
	}

	return append(pwdHash, salt...), nil
}

func verifyPassphrase(passphrase string, targetKey []byte) (bool, error) {
	keyLenBytes := len(targetKey) - common.ScryptSaltLen
	if keyLenBytes < 1 {
		return false, errors.New("Invalid targetKey length")
	}

	targetMasterKey := targetKey[:keyLenBytes]
	salt := targetKey[keyLenBytes:]

	sourceMasterKey, err := scrypt.Key([]byte(passphrase), salt, common.ScryptN, common.ScryptR, common.ScryptP, common.ScryptKeyLen)
	if err != nil {
		return false, fmt.Errorf("VerifyPassphrase scrypt.Key error: %w", err)
	}

	return bytes.Equal(sourceMasterKey, targetMasterKey), nil
}

func (s *Service) CreateToken(userID string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(common.JWTExpires).Unix()
	td.RtExpires = time.Now().Add(common.JWTRefreshExpires).Unix()

	var err error
	//Creating Access Token
	_, accessTokenSecret, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("CreateToken Keystore.Get error: %w userID %s", err, userID)
	}

	userECDSAKey, err := crypto.ToECDSA(accessTokenSecret[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w userID %s", err, userID)
	}

	atClaims := TokenClaims{}
	atClaims.ExpiresAt = td.AtExpires
	atClaims.TokenType = TokenAccessType

	// TODO to fill user metadata like roles we should create new method in contract i.e. UserGet!!!
	at := jwt.NewWithClaims(jwt.SigningMethodES256, atClaims)

	td.AccessToken, err = at.SignedString(userECDSAKey)
	if err != nil {
		return nil, fmt.Errorf("at.SignedString error:%w", err)
	}

	rtClaims := TokenClaims{}
	rtClaims.ExpiresAt = td.RtExpires
	rtClaims.TokenType = TokenRefreshType
	rt := jwt.NewWithClaims(jwt.SigningMethodES256, rtClaims)
	td.RefreshToken, err = rt.SignedString(userECDSAKey)

	if err != nil {
		return nil, fmt.Errorf("rt.SignedString error:%w", err)
	}

	return td, nil
}

func (s *Service) ExtractToken(bearToken string) string {
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (s *Service) VerifyAccess(userID, tokenString string) error {
	tokenString = s.ExtractToken(tokenString)

	tokenAccess, err := s.VerifyToken(userID, tokenString, TokenAccessType)
	if err != nil {
		return fmt.Errorf("VerifyToken error: %w", err)
	}

	_, err = s.ExtractTokenMetadata(tokenAccess)
	if err != nil {
		return fmt.Errorf("ExtractTokenMetadata error: %w", err)
	}

	return nil
}

func (s *Service) VerifyToken(userID, tokenString string, tokenType TokenType) (*jwt.Token, error) {
	tokenUUID := userID

	_, tokenSecret, err := s.Infra.Keystore.Get(tokenUUID)
	if err != nil {
		return nil, fmt.Errorf("VerifyToken Keystore.Get error: %w userID %s", err, userID)
	}

	userECDSAKey, err := crypto.ToECDSA(tokenSecret[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w userID %s", err, userID)
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("%w signing method: %v", errors.ErrIsUnsupported, token.Header["alg"])
		}
		return userECDSAKey.Public(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("VerifyToken jwt.Parse error: %w", err)
	}

	c, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrIsNotValid
	}

	if c.TokenType != tokenType {
		return nil, errors.ErrIsNotValid
	}

	if s.IsTokenInBlackList(tokenString) {
		return nil, errors.ErrIsNotValid
	}

	return token, nil
}

func (s *Service) ExtractTokenMetadata(token *jwt.Token) (*TokenClaims, error) {
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.ErrIsNotValid
	}

	return claims, nil
}

func (s *Service) IsTokenInBlackList(tokenRaw string) bool {
	hash := s.GetTokenHash(tokenRaw)
	_, ok := s.Cache.Get(hash)
	return ok
}

func (s *Service) AddTokenInBlackList(tokenRaw string, expires int64) {
	hash := s.GetTokenHash(tokenRaw)
	s.Cache.Set(hash, nil, time.Until(time.Unix(expires, 0)))
}

func (s *Service) GetTokenHash(tokenRaw string) [32]byte {
	return sha3.Sum256([]byte(tokenRaw))
}

func (s *Service) VerifyAndGetTokenDetails(userID, accessToken, refreshToken string) (*TokenDetails, error) {
	details := TokenDetails{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AtExpires:    0,
		RtExpires:    0,
	}

	if accessToken != "" {
		tokenAccess, err := s.VerifyToken(userID, accessToken, TokenAccessType)
		if err != nil {
			return nil, errors.ErrAccessTokenExp
		}

		metadataAccess, err := s.ExtractTokenMetadata(tokenAccess)
		if err != nil {
			return nil, errors.ErrUnauthorized
		}

		details.AtExpires = metadataAccess.ExpiresAt
	}

	if refreshToken != "" {
		tokenRefresh, err := s.VerifyToken(userID, refreshToken, TokenRefreshType)
		if err != nil {
			return nil, errors.ErrRefreshTokenExp
		}

		metadataRefresh, err := s.ExtractTokenMetadata(tokenRefresh)
		if err != nil {
			return nil, errors.ErrUnauthorized
		}

		details.RtExpires = metadataRefresh.ExpiresAt
	}

	return &details, nil
}
