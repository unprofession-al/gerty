package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/unprofession-al/gerty/transformers"
)

type WebHook struct {
	JobName  string
	BaseUrl  string
	Token    string
	FileName string
}

func (wh *WebHook) Create(next http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			m := req.Method
			if wh.BaseUrl != "NONE" && (m == "POST" || m == "PUT" || m == "PATCH" || m == "DELETE") {
				message := "User " + req.Header.Get(HeaderUserName) + " " + m + "ed to " + req.RequestURI
				data, err := transformers.NewAnsibleInventory()
				if err != nil {
					panic("WebHook could not be rendered")
				}
				out, err := json.Marshal(data)
				if err != nil {
					panic("WebHook could not be marshelled")
				}

				params := make(map[string]string)
				params["json"] = "{'parameter': [{'name':'" + wh.FileName + "', 'file':'file0'}, {'name':'message', 'value':'" + message + "'}]}"

				url := wh.BaseUrl + "/job/" + wh.JobName + "/build?token=" + wh.Token

				request, err := makeRequest(url, params, wh.FileName, out)
				if err != nil {
					panic("Could not create request")
				}
				client := &http.Client{}
				response, err := client.Do(request)
				if err != nil || response.StatusCode < 200 || response.StatusCode >= 300 {
					panic("Could not post to WebHook")
				}
			}
		}()
		next.ServeHTTP(res, req)
	}

	return http.HandlerFunc(fn)
}

func makeRequest(uri string, params map[string]string, filename string, data []byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file0", filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, bytes.NewReader(data))

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return req, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}
