package HTTPClientMock

import (
	"bytes"
)

type HTTPClientMock struct {
	postRes []byte
}

func New() *HTTPClientMock {
	return &HTTPClientMock{}
}

func (h *HTTPClientMock) PostFile(url string, data *[]byte) (res []byte, err error) {
	res = h.postRes
	return
}

func (h *HTTPClientMock) Post(url string, body *bytes.Buffer, headers map[string]string) (resBody []byte, err error) {
	resBody = h.postRes
	return
}

func (h *HTTPClientMock) SetPostRes(res []byte) {
	h.postRes = res
}
