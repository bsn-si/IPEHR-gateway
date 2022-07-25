package HTTPClient

import "bytes"

type Interface interface {
	Post(url string, body *bytes.Buffer, headers map[string]string) (resBody []byte, err error)
	PostFile(url string, data *[]byte) (cnt []byte, err error)
}
