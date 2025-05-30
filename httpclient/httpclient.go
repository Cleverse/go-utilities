package httpclient

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/Cleverse/go-utilities/logger"
	"github.com/cockroachdb/errors"
	"github.com/valyala/fasthttp"
)

var defaultFasthttpClient = fasthttp.Client{
	MaxConnsPerHost: 10240, // default is 512
	ReadBufferSize:  4 * 1024,
	WriteBufferSize: 4 * 1024,
}

type Config struct {
	// Enable debug mode
	Debug bool

	// Default headers
	Headers map[string]string
}

type Client struct {
	baseURL        *url.URL
	fasthttpClient *fasthttp.Client
	Config
}

func New(baseURL string, config ...Config) (*Client, error) {
	client, err := NewFromClient(&defaultFasthttpClient, baseURL, config...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client, nil
}

func NewFromClient(client *fasthttp.Client, baseURL string, config ...Config) (*Client, error) {
	var parsedBaseURL *url.URL
	var err error
	if baseURL == "" {
		parsedBaseURL = &url.URL{}
	} else {
		parsedBaseURL, err = url.Parse(baseURL)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse base url")
		}
	}
	var cf Config
	if len(config) > 0 {
		cf = config[0]
	}
	if len(cf.Headers) == 0 {
		cf.Headers = make(map[string]string)
	}
	return &Client{
		baseURL:        parsedBaseURL,
		Config:         cf,
		fasthttpClient: client,
	}, nil
}

type RequestOptions struct {
	path     string
	method   string
	Body     []byte
	Query    url.Values
	Header   map[string]string
	FormData url.Values
}

type HttpResponse struct {
	URL string
	fasthttp.Response
}

func (r *HttpResponse) UnmarshalBody(out any) error {
	body, err := r.BodyUncompressed()
	if err != nil {
		return errors.Wrapf(err, "can't uncompress body from %v", r.URL)
	}
	contentType := strings.ToLower(string(r.Header.ContentType()))
	switch {
	case strings.Contains(contentType, "application/json"):
		if err := json.Unmarshal(body, out); err != nil {
			return errors.Wrapf(err, "can't unmarshal json body from %s, %q", r.URL, string(body))
		}
		return nil
	case strings.Contains(contentType, "text/plain"):
		return errors.Errorf("can't unmarshal plain text %q", string(body))
	default:
		return errors.Errorf("unsupported content type: %s, contents: %v", r.Header.ContentType(), string(r.Body()))
	}
}

func (h *Client) request(ctx context.Context, reqOptions RequestOptions) (*HttpResponse, error) {
	start := time.Now()
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(reqOptions.method)
	for k, v := range h.Headers {
		req.Header.Set(k, v)
	}
	for k, v := range reqOptions.Header {
		req.Header.Set(k, v)
	}

	baseUrl := h.BaseURL()
	baseUrl.Path = path.Join(baseUrl.Path, reqOptions.path)
	// Because path.Join cleans the joined path. If path ends with /, append "/" to parsedUrl.Path
	if strings.HasSuffix(reqOptions.path, "/") && !strings.HasSuffix(baseUrl.Path, "/") {
		baseUrl.Path += "/"
	}
	baseQuery := baseUrl.Query()
	for k, v := range reqOptions.Query {
		baseQuery[k] = v
	}
	baseUrl.RawQuery = baseQuery.Encode()

	// remove %20 from requestUrl (empty space)
	requestUrl := strings.TrimSuffix(baseUrl.String(), "%20")
	requestUrl = strings.Replace(requestUrl, "%20?", "?", 1)

	// validate requestUrl
	if _, err := url.Parse(requestUrl); err != nil {
		return nil, errors.Wrapf(err, "can't parse request url: %s", requestUrl)
	}

	req.SetRequestURI(requestUrl)
	if reqOptions.Body != nil {
		req.Header.SetContentType("application/json")
		req.SetBody(reqOptions.Body)
	} else if reqOptions.FormData != nil {
		req.Header.SetContentType("application/x-www-form-urlencoded")
		req.SetBodyString(reqOptions.FormData.Encode())
	}

	resp := fasthttp.AcquireResponse()
	startDo := time.Now()

	defer func() {
		if h.Debug {
			ctx := logger.WithContext(ctx,
				slog.String("method", reqOptions.method),
				slog.String("url", requestUrl),
				slog.Duration("duration", time.Since(start)),
				slog.Duration("latency", time.Since(startDo)),
				slog.Int("req_header_size", len(req.Header.Header())),
				slog.Int("req_content_length", req.Header.ContentLength()),
			)

			if resp.StatusCode() >= 0 {
				ctx = logger.WithContext(ctx,
					slog.Int("status_code", resp.StatusCode()),
					slog.String("resp_content_type", string(resp.Header.ContentType())),
					slog.String("resp_content_encoding", string(resp.Header.ContentEncoding())),
					slog.Int("resp_content_length", len(resp.Body())),
				)
			}

			logger.InfoContext(ctx, "Finished make request", slog.String("package", "httpclient"))
		}

		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	resultCh := make(chan error, 1)

	go func() {
		resultCh <- errors.WithStack(h.fasthttpClient.Do(req, resp))
	}()

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-resultCh:
	}
	if err != nil {
		return nil, errors.Wrapf(err, "error during request: url: %s", requestUrl)
	}

	httpResponse := HttpResponse{
		URL: requestUrl,
	}
	resp.CopyTo(&httpResponse.Response)

	return &httpResponse, nil
}

// BaseURL returns the cloned base URL of the client.
func (h *Client) BaseURL() *url.URL {
	u := *h.baseURL
	return &u
}

func (h *Client) Do(ctx context.Context, method, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	reqOptions.path = path
	reqOptions.method = method
	return h.request(ctx, reqOptions)
}

func (h *Client) Get(ctx context.Context, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	reqOptions.path = path
	reqOptions.method = fasthttp.MethodGet
	return h.request(ctx, reqOptions)
}

func (h *Client) Post(ctx context.Context, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	reqOptions.path = path
	reqOptions.method = fasthttp.MethodPost
	return h.request(ctx, reqOptions)
}

func (h *Client) Put(ctx context.Context, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	reqOptions.path = path
	reqOptions.method = fasthttp.MethodPut
	return h.request(ctx, reqOptions)
}

func (h *Client) Patch(ctx context.Context, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	reqOptions.path = path
	reqOptions.method = fasthttp.MethodPatch
	return h.request(ctx, reqOptions)
}

func (h *Client) Delete(ctx context.Context, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	reqOptions.path = path
	reqOptions.method = fasthttp.MethodDelete
	return h.request(ctx, reqOptions)
}

// Do is a shortcut for New(path).Do(ctx, method, "", reqOptions)
func Do(ctx context.Context, method, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	client, err := New(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client.Do(ctx, method, "", reqOptions)
}

// Get is a shortcut for New(path).Get(ctx, path, reqOptions)
func Get(ctx context.Context, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	client, err := New(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client.Get(ctx, path, reqOptions)
}

// Post is a shortcut for New(path).Post(ctx, path, reqOptions)
func Post(ctx context.Context, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	client, err := New(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client.Post(ctx, path, reqOptions)
}

// Put is a shortcut for New(path).Put(ctx, path, reqOptions)
func Put(ctx context.Context, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	client, err := New(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client.Put(ctx, path, reqOptions)
}

// Patch is a shortcut for New(path).Patch(ctx, path, reqOptions)
func Patch(ctx context.Context, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	client, err := New(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client.Patch(ctx, path, reqOptions)
}

// Delete is a shortcut for New(path).Delete(ctx, path, reqOptions)
func Delete(ctx context.Context, path string, reqOptions RequestOptions) (*HttpResponse, error) {
	client, err := New(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client.Delete(ctx, path, reqOptions)
}
