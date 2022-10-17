package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/akyoto/cache"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"

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

type TokenMetadata struct {
	UUID string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

const (
	N      = 1048576
	r      = 8
	p      = 1
	keyLen = 32
)

func NewUserService(cfg *config.Config, infra *infrastructure.Infra) *Service {
	p := processing.New(
		infra.LocalDB,
		infra.EthClient,
		infra.FilecoinClient,
		infra.IpfsClient,
		cfg.Storage.Localfile.Path,
	)

	p.Start()

	return &Service{
		Infra: infra,
		Proc:  p,
		Cache: cache.New(common.CacheCleanerTimeout),
	}
}

func (s *Service) Register(ctx context.Context, procRequest *proc.Request, user *model.UserCreateRequest) (err error) {
	ehrSystemID := ctx.(*gin.Context).GetString("ehrSystemID")
	address, pwdHash, err := s.getUserAddressAndHash(ehrSystemID, user.UserID, user.Password)

	if err != nil {
		return fmt.Errorf("getUserAddressAndHash error: %w", err)
	}

	requestID := ctx.(*gin.Context).GetString("reqId")

	txHash, err := s.Infra.Index.UserAdd(requestID, address, user.UserID, user.Role, pwdHash)
	if err != nil {
		return fmt.Errorf("Index.UserAdd error: %w", err)
	}

	procRequest.AddEthereumTx(proc.TxDeleteDoc, txHash)

	return nil
}

func (s *Service) Login(ctx context.Context, user *model.UserAuthRequest) (err error) {
	ehrSystemID := ctx.(*gin.Context).GetString("ehrSystemID")
	address, pwdHash, err := s.getUserAddressAndHash(ehrSystemID, user.UserID, user.Password)

	if err != nil {
		return fmt.Errorf("Login s.getUserAddressAndHash error: %w", err)
	}

	userHash, err := s.Infra.Index.GetUserPasswordHash(ctx, address)
	if err != nil {
		return fmt.Errorf("Login.GetUserPasswordHash error: %w", err)
	}

	if !bytes.Equal(pwdHash, userHash) {
		return errors.ErrFieldIsIncorrect("Password")
	}

	return nil
}

func (s *Service) getUserAddressAndHash(ehrSystemID, userID, password string) (eth_common.Address, []byte, error) {
	_, userPrivateKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return eth_common.Address{}, nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	privateUserKey := userPrivateKey[:]

	privateKey, err := crypto.ToECDSA(privateUserKey)
	if err != nil {
		return eth_common.Address{}, nil, fmt.Errorf("crypto.ToECDSA error: %w userID %s", err, userID)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	pwdHash, err := s.generateHash(privateUserKey, userID, ehrSystemID, password)
	if err != nil {
		return eth_common.Address{}, nil, fmt.Errorf("register s.generateHash error: %w userID %s, password: %s", err, userID, password)
	}
	return address, pwdHash, nil
}

// method should be idempotent
func (s *Service) generateHash(salt []byte, phrases ...string) ([]byte, error) {
	password := strings.Join(phrases, "")
	// todo salt - random
	pwdHash, err := scrypt.Key([]byte(password), salt, N, r, p, keyLen)
	if err != nil {
		return nil, fmt.Errorf("scrypt.Key error: %w", err)
	}

	return pwdHash, nil
}

func (s *Service) CreateToken(userID string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(common.JWTExpires).Unix()
	td.AccessUUID = uuid.New().String()

	td.RtExpires = time.Now().Add(common.JWTRefreshExpires).Unix()
	td.RefreshUUID = td.AccessUUID + "++" + userID

	var err error
	//Creating Access Token
	accessTokenSecret, _, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("CreateToken Keystore.Get error: %w userID %s", err, userID)
	}

	refreshTokenSecret, _, err := s.Infra.Keystore.Get(userID + "_refresh")
	if err != nil {
		return nil, fmt.Errorf("CreateRefreshToken Keystore.Get error: %w userID %s", err, userID)
	}

	atClaims := jwt.MapClaims{}
	atClaims["uuid"] = td.AccessUUID
	atClaims["exp"] = td.AtExpires
	// TODO to fill user metadata like roles we should create new method in contract i.e. UserGet!!!
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString((*accessTokenSecret)[:])

	if err != nil {
		return nil, fmt.Errorf("at.SignedString error:%w", err)
	}

	//Creating Refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["uuid"] = td.RefreshUUID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString((*refreshTokenSecret)[:])

	if err != nil {
		return nil, fmt.Errorf("rt.SignedString error:%w", err)
	}

	return td, nil
}

func (s *Service) CreateAuth(userid string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	s.Cache.Set(td.AccessUUID, userid, at.Sub(now))
	s.Cache.Set(td.RefreshUUID, userid, rt.Sub(now))

	return nil
}

func (s *Service) ExtractToken(bearToken string) string {
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (s *Service) VerifyToken(userID, tokenString string, isRefreshToken bool) (*jwt.Token, error) {
	tokenUUID := userID
	if isRefreshToken {
		tokenUUID = tokenUUID + "_refresh"
	}

	tokenSecret, _, err := s.Infra.Keystore.Get(tokenUUID)
	if err != nil {
		return nil, fmt.Errorf("VerifyToken Keystore.Get error: %w userID %s", err, userID)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w signing method: %v", errors.ErrIsUnsupported, token.Header["alg"])
		}
		return (*tokenSecret)[:], nil
	})

	if err != nil {
		return nil, fmt.Errorf("VerifyToken jwt.Parse error: %w", err)
	}

	//Since token is valid, get the uuid:
	_, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if !ok || !token.Valid {
		return nil, errors.ErrIsNotValid
	}

	return token, nil
}

func (s *Service) ExtractTokenMetadata(token *jwt.Token) (*TokenMetadata, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrIsNotValid
	}

	u, ok := claims["uuid"].(string)
	if !ok {
		return nil, errors.ErrIsEmpty
	}

	return &TokenMetadata{
		UUID: u,
	}, nil
}

func (s *Service) DeleteToken(token string) {
	s.Cache.Delete(token)

	return
}
