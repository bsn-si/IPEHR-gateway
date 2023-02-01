package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/akyoto/cache"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/compressor"
	docModel "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	proc "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/users"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/infrastructure"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/roles"
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

func NewService(infra *infrastructure.Infra, p *processing.Proc) *Service {
	return &Service{
		Infra: infra,
		Cache: cache.New(common.CacheCleanerTimeout),
		Proc:  p,
	}
}

func (s *Service) NewProcRequest(reqID, userID string, kind processing.RequestKind) (processing.RequestInterface, error) {
	return s.Proc.NewRequest(reqID, userID, "", kind)
}

func (s *Service) Register(ctx context.Context, user *model.UserCreateRequest, systemID, reqID string) (err error) {
	_, userPrivKey, err := s.Infra.Keystore.Get(user.UserID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, user.UserID)
	}

	pwdHash, err := generateHashFromPassword(systemID, user.UserID, user.Password)
	if err != nil {
		return fmt.Errorf("generateHashFromPassword error: %w", err)
	}

	var content []byte

	switch roles.Role(user.Role) {
	case roles.Patient:
	case roles.Doctor:
		info := model.UserInfo{
			UserID:      user.UserID,
			Name:        user.Name,
			Address:     user.Address,
			Description: user.Description,
			PictuteURL:  user.PictuteURL,
		}

		content, err = msgpack.Marshal(info)
		if err != nil {
			return fmt.Errorf("msgpack.Marshal error: %w", err)
		}

		content, err = compressor.New(compressor.BestCompression).Compress(content)
		if err != nil {
			return fmt.Errorf("UserInfo content compression error: %w", err)
		}
	default:
		return errors.ErrFieldIsIncorrect("user.Role")
	}

	procRequest, err := s.NewProcRequest(reqID, user.UserID, processing.RequestUserRegister)
	if err != nil {
		return fmt.Errorf("NewProcRequest error: %w", err)
	}

	multiCallTx, err := s.Infra.Index.MultiCallUsersNew(ctx, userPrivKey)
	if err != nil {
		return fmt.Errorf("MultiCallUsersNew error: %w. userID: %s", err, user.UserID)
	}

	userNewPacked, err := s.Infra.Index.UserNew(ctx, user.UserID, systemID, user.Role, pwdHash, content, userPrivKey, nil)
	if err != nil {
		return fmt.Errorf("Index.UserNew error: %w", err)
	}

	multiCallTx.Add(uint8(proc.TxUserNew), userNewPacked)

	if user.Role == uint8(roles.Patient) {
		// 'doctors' userGroup creating
		{
			groupName := common.DefaultGroupDoctors
			groupDescription := ""

			userGroupCreatePacked, _, err := s.groupCreatePack(ctx, user.UserID, groupName, groupDescription, nil)
			if err != nil {
				return fmt.Errorf("service.GroupCreate error: %w", err)
			}

			multiCallTx.Add(uint8(proc.TxUserGroupCreate), userGroupCreatePacked)
		}
	}

	txHash, err := multiCallTx.Commit()
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return errors.ErrNotFound
		} else if strings.Contains(err.Error(), "AEX") {
			return errors.ErrAlreadyExist
		}

		return fmt.Errorf("UserRegister multicall commit error: %w", err)
	}

	for _, txKind := range multiCallTx.GetTxKinds() {
		procRequest.AddEthereumTx(proc.TxKind(txKind), txHash)
	}

	if err := procRequest.Commit(); err != nil {
		return fmt.Errorf("User register procRequest commit error: %w", err)
	}

	return nil
}

func (s *Service) Login(ctx context.Context, userID, systemID, password string) (err error) {
	address, err := s.getUserAddress(userID)
	if err != nil {
		return fmt.Errorf("Login s.getUserAddress error: %w", err)
	}

	pwdHash, err := s.Infra.Index.GetUserPasswordHash(ctx, address)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return err
		}
		return fmt.Errorf("Login.GetUserPasswordHash error: %w", err)
	}

	match, err := verifyPassphrase(userID+systemID+password, pwdHash)
	if err != nil {
		return fmt.Errorf("verifyPassphrase error: %w", err)
	}

	if !match {
		return errors.ErrAuthorization
	}

	return nil
}

func (s *Service) Info(ctx context.Context, userID, systemID string) (*model.UserInfo, error) {
	address, err := s.getUserAddress(userID)
	if err != nil {
		return nil, fmt.Errorf("s.getUserAddress error: %w", err)
	}

	user, err := s.Infra.Index.GetUser(ctx, address)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("Login.GetUserPasswordHash error: %w", err)
	}

	userInfo, err := extractUserInfo(user)
	if err != nil {
		return nil, fmt.Errorf("extractUserInfo error: %w", err)
	}

	if userInfo.Role == roles.Patient.String() {
		ehrID, err := s.Infra.Index.GetEhrUUIDByUserID(ctx, userID, systemID)
		if err != nil {
			return nil, fmt.Errorf("Info.GetEhrUUIDByID error: %w", err)
		}

		userInfo.EhrID = ehrID
	}

	return userInfo, nil
}

func (s *Service) InfoByCode(ctx context.Context, code int) (*model.UserInfo, error) {
	user, err := s.Infra.Index.GetUserByCode(ctx, uint64(code))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("Login.GetUserPasswordHash error: %w", err)
	}

	userInfo, err := extractUserInfo(user)
	if err != nil {
		return nil, fmt.Errorf("extractUserInfo error: %w", err)
	}

	return userInfo, nil
}

func extractUserInfo(user *users.IUsersUser) (*model.UserInfo, error) {
	var (
		userInfo model.UserInfo
		err      error
	)

	if user.Role == uint8(roles.Doctor) {
		content := docModel.AttributesUsers(user.Attrs).GetByCode(docModel.AttributeContent)
		if content == nil {
			return nil, errors.ErrFieldIsEmpty("AttributeContent")
		}

		content, err = compressor.New(compressor.BestCompression).Decompress(content)
		if err != nil {
			return nil, fmt.Errorf("DoctorInfo content decompression error: %w", err)
		}

		err = msgpack.Unmarshal(content, &userInfo)
		if err != nil {
			return nil, fmt.Errorf("msgpack.Marshal error: %w", err)
		}

		codeInt := binary.BigEndian.Uint64(user.IDHash[0:8]) % common.UserCodeMask
		userInfo.Code = fmt.Sprintf("%08d", codeInt)
	}

	timestamp := docModel.AttributesUsers(user.Attrs).GetByCode(docModel.AttributeTimestamp)
	if timestamp == nil {
		return nil, errors.ErrFieldIsEmpty("AttributeTimestamp")
	}

	userInfo.Role = roles.Role(user.Role).String()

	userInfo.TimeCreated = time.Unix(big.NewInt(0).SetBytes(timestamp).Int64(), 0).Format(common.OpenEhrTimeFormat)

	return &userInfo, nil
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
