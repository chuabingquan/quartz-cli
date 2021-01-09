package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/mholt/archiver"
)

const (
	endpoint = "http://localhost:8080/api/v0/jobs"
)

// Model ...
type Model struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Cron ...
type Cron struct {
	ID         string `json:"id"`
	Expression string `json:"pattern"`
}

// Job ...
type Job struct {
	Model
	ID       string `json:"id"`
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
	Schedule []Cron `json:"schedule"`
}

// JobService ...
type JobService struct{}

// NewJobService ...
func NewJobService() JobService {
	return JobService{}
}

// Jobs ...
func (js *JobService) Jobs() ([]Job, error) {
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	jobs := []Job{}
	err = json.Unmarshal(data, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// Job ...
func (js *JobService) Job(jobID string) {

}

// CreateJob ...
func (js *JobService) CreateJob(src string) (string, error) {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return "", err
	}

	fileNames := []string{}
	for _, fileName := range files {
		fileNames = append(fileNames, src+"/"+fileName.Name())
	}

	err = archiver.Archive(fileNames, src+"/app.tar.gz")
	if err != nil {
		return "", err
	}

	extraParams := map[string]string{}
	request, err := newfileUploadRequest(endpoint, extraParams, "file", src+"/app.tar.gz")
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Deployment failed for unexpected reasons")
	}

	err = os.Remove(src + "/app.tar.gz")
	if err != nil {
		return "", err
	}

	return "Project is successfully deployed.", nil
}

// DeleteJob ...
func (js *JobService) DeleteJob(jobID string) (string, error) {
	url, _ := url.Parse(endpoint + "/" + jobID)
	req := &http.Request{
		Method: http.MethodDelete,
		URL:    url,
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Failed to teardown project")
	}

	return "Project is successfully torn down.", nil
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
