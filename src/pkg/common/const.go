package common

import "time"

const (
	OpenEhrTimeFormat                = "2006-01-02T15:04:05.999-07:00"
	EhrSystemID                      = "openEHRSys.example.com"
	PageLimit                        = 10
	BlockchainTxProcAwaitTime        = time.Second * 1
	FilecoinTxProcAwaitTime          = time.Second * 5
	JWTExpires                       = time.Minute * 15
	JWTRefreshExpires                = time.Hour * 24 * 7
	CacheCleanerTimeout              = 5 * time.Minute
	RegisterRequestTimeout           = time.Second * 60
	UserCodeMask              uint64 = 99999999 // (8 digits)
	DefaultGroupAllDocuments         = "all documents"
	DefaultGroupDoctors              = "doctors"
	QueryExecutionTimeout            = time.Second * 60
	WebRequestTimeout                = time.Second * 60

	ScryptKeyLen  = 32
	ScryptSaltLen = 16
	ScryptN       = 1048576
	ScryptR       = 8
	ScryptP       = 1
)
