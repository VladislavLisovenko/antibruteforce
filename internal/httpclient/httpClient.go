package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Get(ctx context.Context, endpoint string, params url.Values) ([]byte, error)
	Post(ctx context.Context, endpoint string, params url.Values) ([]byte, error)
	Delete(ctx context.Context, endpoint string, params url.Values) ([]byte, error)
}

type httpClient struct {
	Host string
	cl   *http.Client
}

func New(host string) HTTPClient {
	cl := http.DefaultClient

	return &httpClient{
		Host: host,
		cl:   cl,
	}
}

func (h *httpClient) Get(ctx context.Context, endpoint string, params url.Values) ([]byte, error) {
	return h.sendRequest(ctx, h.buildURL(endpoint, params), http.MethodGet)
}

func (h *httpClient) Post(ctx context.Context, endpoint string, params url.Values) ([]byte, error) {
	return h.sendRequest(ctx, h.buildURL(endpoint, params), http.MethodPost)
}

func (h *httpClient) Delete(ctx context.Context, endpoint string, params url.Values) ([]byte, error) {
	return h.sendRequest(ctx, h.buildURL(endpoint, params), http.MethodDelete)
}

func (h *httpClient) sendRequest(ctx context.Context, uri string, method string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := h.cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body read error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code %d, body: %s", resp.StatusCode, string(body))
	}

	return body, err
}

func (h *httpClient) buildURL(endpoint string, params url.Values) string {
	return h.Host + "/" + endpoint + "?" + params.Encode()
}
