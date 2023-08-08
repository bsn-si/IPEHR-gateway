package processing

type RequestKind uint8

const (
	RequestUnknown RequestKind = iota
	RequestEhrCreate
	RequestEhrCreateWithID
	RequestEhrGetBySubject
	RequestEhrGetByID
	RequestEhrStatusCreate
	RequestEhrStatusUpdate
	RequestEhrStatusGetByID
	RequestEhrStatusGetByTime
	RequestCompositionCreate
	RequestCompositionUpdate
	RequestCompositionGetByID
	RequestCompositionDelete
	RequestUserRegister
	RequestDocAccessSet
	RequestQueryStore
	RequestUserGroupCreate
	RequestUserGroupAddUser
	RequestUserGroupRemoveUser
	RequestContributionCreate
	RequestTemplateCreate
	RequestDirectoryCreate
	RequestDirectoryUpdate
	RequestDirectoryDelete
)

var reqKinds = map[RequestKind]string{
	RequestEhrCreate:          "EhrCreate",
	RequestEhrGetBySubject:    "EhrGetBySubject",
	RequestEhrGetByID:         "EhrGetByID",
	RequestEhrStatusCreate:    "EhrStatusCreate",
	RequestEhrStatusUpdate:    "EhrStatusUpdate",
	RequestEhrStatusGetByID:   "EhrStatusGetByID",
	RequestEhrStatusGetByTime: "EhrStatusGetByTime",
	RequestCompositionCreate:  "CompositionCreate",
	RequestCompositionUpdate:  "CompositionUpdate",
	RequestCompositionGetByID: "CompositionGetByID",
	RequestCompositionDelete:  "CompositionDelete",
	RequestTemplateCreate:     "TemplateStore",
	RequestDirectoryCreate:    "DirectoryCreate",
	RequestDirectoryUpdate:    "DirectoryUpdate",
	RequestDirectoryDelete:    "DirectoryDelete",
}

func (k RequestKind) String() string {
	if rk, ok := reqKinds[k]; ok {
		return rk
	}

	return reqKinds[RequestUnknown]
}
