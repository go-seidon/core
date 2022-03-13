package http

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Header struct {
	Key, Value string
}

type RequestParam struct {
	Url     string
	Headers []Header
	Timeout time.Duration
	Body    io.Reader
}

type ResponseResult struct {
	Data       []byte
	StatusCode int
}

type HttpClient interface {
	Get(RequestParam) (*ResponseResult, error)
	Post(RequestParam) (*ResponseResult, error)
}

type httpService struct {
	httpClient http.Client
}

func (s *httpService) Get(p RequestParam) (*ResponseResult, error) {
	req, err := http.NewRequest("GET", p.Url, p.Body)
	if err != nil {
		return nil, err
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := &ResponseResult{
		StatusCode: res.StatusCode,
		Data:       data,
	}

	return result, nil
}

func (s *httpService) Post(p RequestParam) (*ResponseResult, error) {
	req, err := http.NewRequest("POST", p.Url, p.Body)
	if err != nil {
		return nil, err
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := &ResponseResult{
		StatusCode: res.StatusCode,
		Data:       data,
	}

	return result, nil
}

func NewHttpClient() *httpService {
	return &httpService{
		httpClient: http.Client{},
	}
}
