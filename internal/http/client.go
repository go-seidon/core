package http

import (
	"io"
	"net/http"
	"time"

	"golang.org/x/net/context"
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

type HttpService interface {
	Get(p RequestParam) (*ResponseResult, error)
	Post(p RequestParam) (*ResponseResult, error)
}

type httpService struct {
	httpClient http.Client
}

func (s *httpService) Get(p RequestParam) (*ResponseResult, error) {
	res, err := s.request("GET", p)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
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
	res, err := s.request("POST", p)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := &ResponseResult{
		StatusCode: res.StatusCode,
		Data:       data,
	}

	return result, nil
}

func (s *httpService) request(m string, p RequestParam) (*http.Response, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, m, p.Url, p.Body)
	if err != nil {
		return nil, err
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewHttpClient() HttpService {
	return &httpService{
		httpClient: http.Client{},
	}
}
