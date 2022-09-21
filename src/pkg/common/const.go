package common

import "time"

const (
	OpenEhrTimeFormat         = "2006-01-02T15:04:05.999-07:00"
	EhrSystemID               = "openEHRSys.example.com"
	PageLimit                 = 10
	BlockchainTxProcAwaitTime = time.Millisecond * 500
	FilecoinTxProcAwaitTime   = time.Second * 5
)
