package HTTPClient

import (
	"bytes"
	"hms/gateway/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
)

type HTTPClient struct{}

func New() *HTTPClient {
	return &HTTPClient{}
}

// PostFile Do a POST request with file
func (h *HTTPClient) PostFile(url string, fileContent *[]byte) (res []byte, err error) {
	body, boundary, err := h.createFileBody(fileContent)
	if err != nil {
		return
	}

	headers := map[string]string{
		"Content-Type": "multipart/form-data; boundary=" + boundary,
	}

	res, err = h.Post(url, body, headers)

	return
}

// Post Do a POST request
func (h *HTTPClient) Post(url string, body *bytes.Buffer, headers map[string]string) (resBody []byte, err error) {
	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}

	for i, v := range headers {
		r.Header.Add(i, v)
	}

	client := &http.Client{}

	resp, err := client.Do(r)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	resBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.RequestError(resp.Status, string(resBody))
	}

	return
}

// Prepare post body with file content
func (h *HTTPClient) createFileBody(data *[]byte) (body *bytes.Buffer, boundary string, err error) {
	buf := bytes.NewBuffer(*data)

	body = &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "file.txt")
	if err != nil {
		return
	}

	_, err = io.Copy(part, buf)
	if err != nil {
		return
	}

	boundary = writer.Boundary()

	err = writer.Close()
	if err != nil {
		return
	}

	return
}
