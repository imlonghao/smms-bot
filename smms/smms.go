package smms

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

var (
	endpoint = "https://sm.ms/api/v2/"
)

type Client struct {
	Token string
}

func (c *Client) Upload(photo io.Reader, filename string) (*UploadJSON, error) {
	var err error
	var writer io.Writer
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if writer, err = w.CreateFormFile("smfile", filename); err != nil {
		return &UploadJSON{}, err
	}
	if _, err = io.Copy(writer, photo); err != nil {
		return &UploadJSON{}, err
	}
	if err = w.Close(); err != nil {
		return &UploadJSON{}, err
	}
	req, err := http.NewRequest("POST", endpoint+"upload", &b)
	if err != nil {
		return &UploadJSON{}, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", c.Token)
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return &UploadJSON{}, err
	}
	defer response.Body.Close()
	var result UploadJSON
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&result); err != nil {
		return &UploadJSON{}, err
	}
	return &result, nil
}

func (c *Client) History() (*HistoryJSON, error) {
	var resp *http.Response
	var err error
	if resp, err = http.Get(endpoint+"history") ; err != nil {
		return &HistoryJSON{}, err
	}
	defer resp.Body.Close()
	var result HistoryJSON
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&result); err != nil {
		return &HistoryJSON{}, err
	}
	return &result, nil
}

func (c *Client) Delete(hash string) (*DeleteJSON, error){
	var resp *http.Response
	var err error
	if resp, err = http.Get(endpoint+"delete/"+hash) ; err != nil {
		return &DeleteJSON{}, err
	}
	defer resp.Body.Close()
	var result DeleteJSON
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&result); err != nil {
		return &DeleteJSON{}, err
	}
	return &result, nil
}

func (c *Client) Clear() (*ClearJSON, error) {
	var resp *http.Response
	var err error
	if resp, err = http.Get(endpoint+"clear") ; err != nil {
		return &ClearJSON{}, err
	}
	defer resp.Body.Close()
	var result ClearJSON
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&result); err != nil {
		return &ClearJSON{}, err
	}
	return &result, nil
}