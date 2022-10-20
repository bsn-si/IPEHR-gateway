package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/akyoto/cache"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
	Exp int64
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

const (
	N                = 1048576
	r                = 8
	p                = 1
	keyLen           = 32
	metadataLenBytes = 28
	saltLenBytes     = 16
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
	address, err := s.getUserAddress(user.UserID)

	if err != nil {
		return fmt.Errorf("getUserAddress error: %w", err)
	}

	pwdHash, err := s.generateHashFromPassword(ehrSystemID, user.UserID, user.Password)

	if err != nil {
		return fmt.Errorf("generateHashFromPassword error: %w", err)
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
	address, err := s.getUserAddress(user.UserID)

	if err != nil {
		return fmt.Errorf("Login s.getUserAddress error: %w", err)
	}

	passwordHash, err := s.Infra.Index.GetUserPasswordHash(ctx, address)
	if err != nil {
		return fmt.Errorf("Login.GetUserPasswordHash error: %w", err)
	}

	match, err := s.VerifyPassphrase(ehrSystemID+user.UserID+user.Password, passwordHash)

	if err != nil {
		return fmt.Errorf("generateHashFromPassword error: %w", err)
	}

	if !match {
		return errors.ErrAuthorization
	}

	return nil
}

func (s *Service) VerifyPassphrase(passphrase string, targetKey []byte) (bool, error) {
	keylenBytes := len(targetKey) - metadataLenBytes
	if keylenBytes < 1 {
		return false, errors.New("Invalid targetKey length")
	}
	// Get the master_key
	targetMasterKey := targetKey[:keylenBytes]
	// Get the salt
	salt := targetKey[keylenBytes : keylenBytes+saltLenBytes]
	// Get the params
	var N, r, p int32

	paramsStartIndex := keylenBytes + saltLenBytes

	err := binary.Read(bytes.NewReader(targetKey[paramsStartIndex:paramsStartIndex+4]), // 4 bytes for N
		binary.LittleEndian,
		&N)
	if err != nil {
		return false, fmt.Errorf("VerifyPassphrase binary.Read error: %w", err)
	}

	err = binary.Read(bytes.NewReader(targetKey[paramsStartIndex+4:paramsStartIndex+8]), // 4 bytes for r
		binary.LittleEndian,
		&r)
	if err != nil {
		return false, fmt.Errorf("VerifyPassphrase binary.Read error: %w", err)
	}

	err = binary.Read(bytes.NewReader(targetKey[paramsStartIndex+8:paramsStartIndex+12]), // 4 bytes for p
		binary.LittleEndian,
		&p)
	if err != nil {
		return false, fmt.Errorf("VerifyPassphrase binary.Read error: %w", err)
	}

	sourceMasterKey, err := scrypt.Key([]byte(passphrase),
		salt,
		int(N), // Must be a power of 2 greater than 1
		int(r),
		int(p), // r*p must be < 2^30
		keylenBytes)
	if err != nil {
		return false, fmt.Errorf("VerifyPassphrase scrypt.Key error: %w", err)
	}

	keyComp := subtle.ConstantTimeCompare(sourceMasterKey, targetMasterKey) != 0
	return keyComp, nil
}

func (s *Service) getUserAddress(userID string) (eth_common.Address, error) {
	_, userPrivateKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return eth_common.Address{}, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	privateUserKey := userPrivateKey[:]

	privateKey, err := crypto.ToECDSA(privateUserKey)
	if err != nil {
		return eth_common.Address{}, fmt.Errorf("crypto.ToECDSA error: %w userID %s", err, userID)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return address, nil
}

func (s *Service) generateHashFromPassword(ehrSystemID, userID, password string) ([]byte, error) {
	salt, err := s.generateSalt()
	if err != nil {
		return nil, fmt.Errorf("s.generateSalt error: %w userID %s, password: %s", err, userID, password)
	}

	pwdHash, err := s.generateHash(salt, userID, ehrSystemID, password)
	if err != nil {
		return nil, fmt.Errorf("s.generateHash error: %w userID %s, password: %s", err, userID, password)
	}

	// Appending the salt
	pwdHash = append(pwdHash, salt...)

	// Encoding the params to be stored
	buf := &bytes.Buffer{}
	for _, elem := range [3]int{N, r, p} {
		err = binary.Write(buf, binary.LittleEndian, int32(elem))
		if err != nil {
			return nil, fmt.Errorf("binary.Write error: %w userID %s, password: %s", err, userID, password)
		}
	}

	pwdHash = append(pwdHash, buf.Bytes()...)

	return pwdHash, nil
}

func (s *Service) generateSalt() ([]byte, error) {
	salt := make([]byte, saltLenBytes)

	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generateSalt rand.Read error: %w", err)
	}

	return salt, nil
}

func (s *Service) generateHash(salt []byte, phrases ...string) ([]byte, error) {
	password := strings.Join(phrases, "")

	pwdHash, err := scrypt.Key([]byte(password), salt, N, r, p, keyLen)
	if err != nil {
		return nil, fmt.Errorf("generateHash scrypt.Key error: %w", err)
	}

	return pwdHash, nil
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

	ecdsaAccessKey, err := crypto.ToECDSA(accessTokenSecret[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w userID %s", err, userID)
	}

	atClaims := jwt.MapClaims{}
	atClaims["exp"] = td.AtExpires
	// TODO to fill user metadata like roles we should create new method in contract i.e. UserGet!!!
	at := jwt.NewWithClaims(jwt.SigningMethodES256, atClaims)
	td.AccessToken, err = at.SignedString(ecdsaAccessKey)

	if err != nil {
		return nil, fmt.Errorf("at.SignedString error:%w", err)
	}

	//Creating Refresh token
	_, refreshTokenSecret, err := s.Infra.Keystore.Get(userID + "_refresh")
	if err != nil {
		return nil, fmt.Errorf("CreateRefreshToken Keystore.Get error: %w userID %s", err, userID)
	}

	ecdsaRefreshKey, err := crypto.ToECDSA(refreshTokenSecret[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w userID %s", err, userID)
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodES256, rtClaims)
	td.RefreshToken, err = rt.SignedString(ecdsaRefreshKey)

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

func (s *Service) VerifyToken(userID, tokenString string, isRefreshToken bool) (*jwt.Token, error) {
	tokenUUID := userID
	if isRefreshToken {
		tokenUUID = tokenUUID + "_refresh"
	}

	_, tokenSecret, err := s.Infra.Keystore.Get(tokenUUID)
	if err != nil {
		return nil, fmt.Errorf("VerifyToken Keystore.Get error: %w userID %s", err, userID)
	}

	ecdsaKey, err := crypto.ToECDSA(tokenSecret[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w userID %s", err, userID)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("%w signing method: %v", errors.ErrIsUnsupported, token.Header["alg"])
		}
		return ecdsaKey.Public(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("VerifyToken jwt.Parse error: %w", err)
	}

	//Since token is valid, get the uuid:
	_, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if !ok || !token.Valid {
		return nil, errors.ErrIsNotValid
	}

	if ok := s.IsTokenInBlackList(tokenString); ok {
		return nil, errors.ErrIsNotValid
	}

	return token, nil
}

func (s *Service) ExtractTokenMetadata(token *jwt.Token) (*TokenMetadata, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrIsNotValid
	}

	return &TokenMetadata{
		Exp: int64(claims["exp"].(float64)),
	}, nil
}

func (s *Service) IsTokenInBlackList(tokenRaw string) bool {
	hash := string(s.GetTokenHash(tokenRaw))
	_, ok := s.Cache.Get(hash)
	return ok
}

func (s *Service) AddTokenInBlackList(tokenRaw string, expires int64) {
	at := time.Unix(expires, 0) //converting Unix to UTC(to Time object)
	now := time.Now()

	hash := string(s.GetTokenHash(tokenRaw))

	s.Cache.Set(hash, nil, at.Sub(now))
}

func (s *Service) GetTokenHash(tokenRaw string) []byte {
	h := sha256.New()
	h.Write([]byte(tokenRaw))

	return h.Sum(nil)
}
