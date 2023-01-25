package template

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/compressor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/parser/adl14"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/parser/adl2"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/status"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/types"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/helper"
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

func (s *Service) Store(ctx context.Context, userID, systemID, reqID string, m *model.Template) error {
	userPubKey, userPrivKey, err := s.docSvc.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	if m.Body == nil {
		return errors.ErrFieldIsEmpty("Body")
	}

	docBytes := make([]byte, len(m.Body))
	copy(docBytes, m.Body)

	if s.docSvc.Infra.Compressor != nil {
		docBytes, err = s.docSvc.Infra.Compressor.Compress(docBytes)
		if err != nil {
			return fmt.Errorf("Compress error: %w", err)
		}
	}

	// Document encryption key generation
	key := chachaPoly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(m.UID))
	if err != nil {
		return fmt.Errorf("EncryptWithAuthData error: %w", err)
	}

	// IPFS saving
	CID, err := s.docSvc.Infra.IpfsClient.Add(ctx, docEncrypted)
	if err != nil {
		return fmt.Errorf("IpfsClient.Add error: %w", err)
	}

	// Filecoin saving
	dealCID, minerAddr, err := s.docSvc.Infra.FilecoinClient.StartDeal(ctx, CID, uint64(len(docEncrypted)))
	if err != nil {
		return fmt.Errorf("FilecoinClient.StartDeal error: %w", err)
	}

	docIDEncrypted, err := key.Encrypt([]byte(m.UID))
	if err != nil {
		return fmt.Errorf("key.Encrypt error: %w", err)
	}

	content, err := msgpack.Marshal(m)
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

	procRequest, err := s.docSvc.Proc.NewRequest(reqID, userID, "", processing.RequestTemplateCreate)
	if err != nil {
		return fmt.Errorf("Proc.NewRequest error: %w", err)
	}

	// Add filecoin tx
	procRequest.AddFilecoinTx(processing.TxSaveTemplate, CID.String(), dealCID.String(), minerAddr)

	multiCallTx, err := s.docSvc.Infra.Index.MultiCallEhrNew(ctx, userPrivKey)
	if err != nil {
		return fmt.Errorf("MultiCallEhrNew error: %w userID %s", err, userID)
	}

	// Index Docs
	keyEncr, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
	if err != nil {
		return fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	CIDEncr, err := keybox.SealAnonymous(CID.Bytes(), userPubKey)
	if err != nil {
		return fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	tmplIDHash := sha3.Sum256([]byte(m.TemplateID))

	docMeta := &model.DocumentMeta{
		Status:    uint8(status.ACTIVE),
		Id:        CID.Bytes(),
		Version:   []byte(m.Version),
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

	packed, err := s.docSvc.Infra.Index.AddEhrDoc(ctx, types.Template, docMeta, userPrivKey, nil)
	if err != nil {
		return fmt.Errorf("Index.AddEhrDoc error: %w", err)
	}

	multiCallTx.Add(uint8(processing.TxAddEhrDoc), packed)

	txHash, err := multiCallTx.Commit()
	if err != nil {
		return fmt.Errorf("Create template commit error: %w", err)
	}

	for _, txKind := range multiCallTx.GetTxKinds() {
		procRequest.AddEthereumTx(processing.TxKind(txKind), txHash)
	}

	if err := procRequest.Commit(); err != nil {
		return fmt.Errorf("ProcRequest commit error: %w", err)
	}

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
