package template

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/compressor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/parser/adl14"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/parser/adl2"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	proc "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/status"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/types"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/helper"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/ehrIndexer"

	"github.com/ipfs/go-cid"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"
)

type Service struct {
	helper.Finder
	docSvc  *service.DefaultDocumentService
	parsers map[string]ADLParser
}

type ADLParser interface {
	Version() model.ADLVer
	IsTypeAllowed(t model.ADLType) bool
	Validate([]byte, model.ADLType) bool
	Parse([]byte, model.ADLType) (*model.Template, error)
	ParseWithFill([]byte, model.ADLType) (*model.Template, error)
}

func NewService(docService *service.DefaultDocumentService) *Service {
	opt14 := adl14.NewParser()
	opt2 := adl2.NewParser()

	ps := map[string]ADLParser{
		opt14.Version(): opt14,
		opt2.Version():  opt2,
	}

	return &Service{
		docSvc:  docService,
		parsers: ps,
	}
}

func (s *Service) Parser(version string) (ADLParser, error) {
	p, ok := s.parsers[version]
	if !ok {
		return nil, errors.ErrIsNotExist
	}

	return p, nil
}

func (s *Service) GetByID(ctx context.Context, userID, systemID, templateID string) (*model.Template, error) {
	tmplIDHash := sha3.Sum256([]byte(templateID))

	docMeta, err := s.docSvc.Infra.Index.GetDocLastByBaseID(ctx, userID, systemID, types.Template, &tmplIDHash)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, errors.ErrNotFound
		}

		return nil, fmt.Errorf("GetLastVersionDocIndexByBaseID error: %w userID %s templateID %s", err, userID, templateID)
	}

	CID, err := cid.Parse(docMeta.Id)
	if err != nil {
		return nil, fmt.Errorf("cid.Parse error: %w", err)
	}

	docUIDEncrypted := docMeta.GetAttr(model.AttributeDocUIDEncr)
	if docUIDEncrypted == nil {
		return nil, errors.ErrFieldIsEmpty("DocUIDEncrypted")
	}

	docDecrypted, err := s.docSvc.GetDocFromStorageByID(ctx, userID, systemID, &CID, []byte(templateID), docUIDEncrypted)
	if err != nil && errors.Is(err, errors.ErrIsInProcessing) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w userID %s storageID %s", err, userID, &CID)
	}

	//TODO support ADL20
	p, err := s.Parser(model.VerADL1_4)
	if err != nil {
		return nil, fmt.Errorf("Parser error: %w", err)
	}

	//TODO detect type
	t, err := p.ParseWithFill(docDecrypted, model.ADLTypeXML)
	if err != nil {
		return nil, fmt.Errorf("ParseWithFill error: %w", err)
	}

	return t, nil
}

func (s *Service) Store(ctx context.Context, userID, systemID, reqID string, tmpl *model.Template) error {
	userPubKey, userPrivKey, err := s.docSvc.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	if tmpl.Body == nil {
		return errors.ErrFieldIsEmpty("Body")
	}

	docBytes := make([]byte, len(tmpl.Body))
	copy(docBytes, tmpl.Body)

	if s.docSvc.Infra.Compressor != nil {
		docBytes, err = s.docSvc.Infra.Compressor.Compress(docBytes)
		if err != nil {
			return fmt.Errorf("Compress error: %w", err)
		}
	}

	key := chachaPoly.GenerateKey()

	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(tmpl.UID))
	if err != nil {
		return fmt.Errorf("EncryptWithAuthData error: %w", err)
	}

	CID, err := s.docSvc.Infra.IpfsClient.Add(ctx, docEncrypted)
	if err != nil {
		return fmt.Errorf("IpfsClient.Add error: %w", err)
	}

	dealCID, minerAddr, err := s.docSvc.Infra.FilecoinClient.StartDeal(ctx, CID, uint64(len(docEncrypted)))
	if err != nil {
		return fmt.Errorf("FilecoinClient.StartDeal error: %w", err)
	}

	procRequest, err := s.docSvc.Proc.NewRequest(reqID, userID, "", proc.RequestTemplateCreate)
	if err != nil {
		return fmt.Errorf("Proc.NewRequest error: %w", err)
	}

	procRequest.AddFilecoinTx(proc.TxSaveTemplate, CID.String(), dealCID.String(), minerAddr)

	multiCallTx := s.docSvc.Infra.Index.MultiCallEhrNew()

	err = s.addMetaData(multiCallTx, key, tmpl, CID, dealCID, minerAddr, userPubKey, userPrivKey)
	if err != nil {
		return fmt.Errorf("addDataIndex error: %w", err)
	}

	txHash, err := multiCallTx.Commit()
	if err != nil {
		return fmt.Errorf("Create template commit error: %w", err)
	}

	for _, txKind := range multiCallTx.GetTxKinds() {
		procRequest.AddEthereumTx(proc.TxKind(txKind), txHash)
	}

	err = s.setDocAccess(ctx, procRequest, userID, systemID, CID, key, access.Owner, userPubKey, userPrivKey)
	if err != nil {
		return fmt.Errorf("setDocAccess error: %w", err)
	}

	if err := procRequest.Commit(); err != nil {
		return fmt.Errorf("ProcRequest commit error: %w", err)
	}

	return nil
}

func (s *Service) addMetaData(multiCallTx *indexer.MultiCallTx, key *chachaPoly.Key, tmpl *model.Template, CID, dealCID *cid.Cid, minerAddr string, userPubKey, userPrivKey *[32]byte) error {
	keyEncr, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
	if err != nil {
		return fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	CIDEncr, err := keybox.SealAnonymous(CID.Bytes(), userPubKey)
	if err != nil {
		return fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	docIDEncrypted, err := key.Encrypt([]byte(tmpl.UID))
	if err != nil {
		return fmt.Errorf("key.Encrypt error: %w", err)
	}

	tmplIDHash := sha3.Sum256([]byte(tmpl.TemplateID))

	content, err := msgpack.Marshal(tmpl)
	if err != nil {
		return fmt.Errorf("msgpack.Marshal error: %w", err)
	}

	contentCompresed, err := compressor.New(compressor.BestCompression).Compress(content)
	if err != nil {
		return fmt.Errorf("Template Compress error: %w", err)
	}

	contentEncr, err := key.Encrypt(contentCompresed)
	if err != nil {
		return fmt.Errorf("key.Encrypt content error: %w", err)
	}

	docMeta := &model.DocumentMeta{
		Status:    uint8(status.ACTIVE),
		Id:        CID.Bytes(),
		Version:   []byte(tmpl.Version),
		Timestamp: uint32(time.Now().Unix()),
		IsLast:    true,
		Attrs: []ehrIndexer.AttributesAttribute{
			{Code: model.AttributeIDEncr, Value: CIDEncr},
			{Code: model.AttributeKeyEncr, Value: keyEncr},
			{Code: model.AttributeDocUIDHash, Value: tmplIDHash[:]},
			{Code: model.AttributeDocUIDEncr, Value: docIDEncrypted},
			{Code: model.AttributeDealCid, Value: dealCID.Bytes()},
			{Code: model.AttributeMinerAddress, Value: []byte(minerAddr)},
			{Code: model.AttributeContentEncr, Value: contentEncr},
		},
	}

	packed, err := s.docSvc.Infra.Index.AddEhrDoc(types.Template, docMeta, userPrivKey)
	if err != nil {
		return fmt.Errorf("Index.AddEhrDoc error: %w", err)
	}

	multiCallTx.Add(uint8(proc.TxAddEhrDoc), packed)

	return nil
}

func (s *Service) setDocAccess(ctx context.Context, req proc.RequestInterface, userID, systemID string, CID *cid.Cid, key *chachaPoly.Key, accessLevel access.Level, userPubKey, userPrivKey *[32]byte) error {
	userIDHash := sha3.Sum256([]byte(userID + systemID))
	docIDHash := indexer.Keccak256(CID.Bytes())

	CIDEncr, err := key.Encrypt(CID.Bytes())
	if err != nil {
		return fmt.Errorf("CID encryption error: %w", err)
	}

	keyEncr, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
	if err != nil {
		return fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	accessObj := indexer.AccessObject{
		Kind:    access.Doc,
		IdHash:  *docIDHash,
		IdEncr:  CIDEncr,
		KeyEncr: keyEncr,
		Level:   accessLevel,
	}

	txHash, err := s.docSvc.Infra.Index.SetAccess(ctx, &userIDHash, &accessObj, userPrivKey)
	if err != nil {
		return fmt.Errorf("Index.SetAccess user to template error: %w", err)
	}

	req.AddEthereumTx(proc.TxSetDocAccess, txHash)

	return nil
}

func (s *Service) GetList(ctx context.Context, userID, systemID string) ([]*model.Template, error) {
	userPubKey, userPrivKey, err := s.docSvc.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	list, err := s.docSvc.Infra.Index.ListDocByType(ctx, userID, systemID, types.Template)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.ListDocByType error: %w", err)
	}

	var result = []*model.Template{}

	for i, dm := range list {
		if !dm.IsLast {
			continue
		}

		dm := dm

		key, err := s.docSvc.KeyFromAttribures(&dm, userPubKey, userPrivKey)
		if err != nil {
			return nil, fmt.Errorf("index %d KeyFromAttribures error: %w", i, err)
		}

		content, err := s.docSvc.ContentFromAttributes(&dm, key)
		if err != nil {
			return nil, fmt.Errorf("index %d ContentFromAttributes error: %w", i, err)
		}

		var tmpl model.Template

		err = msgpack.Unmarshal(content, &tmpl)
		if err != nil {
			return nil, fmt.Errorf("index %d Template content unmarshal error: %w", i, err)
		}

		tmpl.CreatedAt = time.Unix(int64(dm.Timestamp), 0).Format(common.OpenEhrTimeFormat)

		result = append(result, &tmpl)
	}

	return result, nil
}

func (s *Service) IsExist(ctx context.Context, args ...string) (bool, error) {
	if len(args) != 3 {
		return false, fmt.Errorf("%w: Expected args: userID, systemID, templateID", errors.ErrCustom)
	}

	userID, systemID, templateID := args[0], args[1], args[2]

	ok, err := s.GetByID(ctx, userID, systemID, templateID)
	if err != nil {
		return false, fmt.Errorf("GetByID error: %w", err)
	}

	return (ok != nil), nil
}
