package processing

type TxKind uint8

const (
	TxUnknown TxKind = iota
	TxMultiCall
	TxSetEhrUser
	TxSetEhrBySubject
	TxSetEhrDocs
	TxSetDocAccess
	TxSetDocGroupAccess
	TxSetUserGroupAccess
	TxDeleteDoc
	TxFilecoinStartDeal
	TxEhrCreateWithID
	TxUpdateEhrStatus
	TxAddEhrDoc
	TxSetDocKeyEncrypted
	TxSaveEhr
	TxSaveEhrStatus
	TxSaveComposition
	TxSaveTemplate
	TxUserRegister
	TxUserNew
	TxUserGroupCreate
	TxUserGroupAddUser
	TxUserGroupRemoveUser
	TxDocGroupCreate
	TxDocGroupAddDoc
	TxIndexDataUpdate
	TxCreateDirectory
)

var txKinds = map[TxKind]string{
	TxMultiCall:           "MultiCall",
	TxSetEhrUser:          "SetEhrUser",
	TxSetEhrBySubject:     "SetEhrBySubject",
	TxSetEhrDocs:          "SetEhrDocs",
	TxSetDocAccess:        "SetDocAccess",
	TxSetDocGroupAccess:   "SetDocGroupAccess",
	TxSetUserGroupAccess:  "SetUserGroupAccess",
	TxDeleteDoc:           "DeleteDoc",
	TxFilecoinStartDeal:   "FilecoinStartDeal",
	TxEhrCreateWithID:     "EhrCreateWithID",
	TxUpdateEhrStatus:     "UpdateEhrStatus",
	TxAddEhrDoc:           "AddEhrDoc",
	TxSetDocKeyEncrypted:  "SetDocKeyEncrypted",
	TxSaveEhr:             "SaveEhr",
	TxSaveEhrStatus:       "SaveEhrStatus",
	TxSaveComposition:     "SaveComposition",
	TxSaveTemplate:        "SaveTemplate",
	TxCreateDirectory:     "CreateDirectory",
	TxUserRegister:        "UserRegister",
	TxUserNew:             "UserNew",
	TxUserGroupCreate:     "UserGroupCreate",
	TxUserGroupAddUser:    "UserGroupAddUser",
	TxUserGroupRemoveUser: "UserGroupRemoveUser",
	TxDocGroupCreate:      "DocGroupCreate",
	TxDocGroupAddDoc:      "DocGroupAddDoc",
	TxIndexDataUpdate:     "IndexDataUpdate",
	TxUnknown:             "Unknown",
}

func (k TxKind) String() string {
	if tk, ok := txKinds[k]; ok {
		return tk
	}

	return txKinds[TxUnknown]
}
