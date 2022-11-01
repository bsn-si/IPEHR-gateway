package model

type DocAccessManageRequest struct {
	ToUserID    string
	CID         string
	AccessLevel uint8
}
