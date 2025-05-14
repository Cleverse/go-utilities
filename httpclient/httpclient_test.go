package httpclient

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

type testResponse struct {
	Message string `json:"message"`
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		config  []Config
		wantErr bool
	}{
		{
			name:    "valid base URL",
			baseURL: "http://example.com",
			config:  []Config{},
			wantErr: false,
		},
		{
			name:    "invalid base URL",
			baseURL: "://invalid",
			config:  []Config{},
			wantErr: true,
		},
		{
			name:    "empty base URL",
			baseURL: "",
			config:  []Config{},
			wantErr: false,
		},
		{
			name:    "with headers",
			baseURL: "http://example.com",
			config: []Config{
				{
					Headers: map[string]string{
						"X-Test": "test-value",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.baseURL, tt.config...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, client)
		})
	}
}

func TestHTTPMethods(t *testing.T) {
	server := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			response := testResponse{Message: string(ctx.Method())}
			ctx.SetContentType("application/json")
			json.NewEncoder(ctx).Encode(response)
		},
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	go server.Serve(ln)

	client, err := New("http://" + ln.Addr().String())
	require.NoError(t, err)

	tests := []struct {
		name     string
		method   func(context.Context, string, RequestOptions) (*HttpResponse, error)
		expected string
	}{
		{
			name: "GET",
			method: func(ctx context.Context, path string, opts RequestOptions) (*HttpResponse, error) {
				return client.Get(ctx, path, opts)
			},
			expected: "GET",
		},
		{
			name: "POST",
			method: func(ctx context.Context, path string, opts RequestOptions) (*HttpResponse, error) {
				return client.Post(ctx, path, opts)
			},
			expected: "POST",
		},
		{
			name: "PUT",
			method: func(ctx context.Context, path string, opts RequestOptions) (*HttpResponse, error) {
				return client.Put(ctx, path, opts)
			},
			expected: "PUT",
		},
		{
			name: "PATCH",
			method: func(ctx context.Context, path string, opts RequestOptions) (*HttpResponse, error) {
				return client.Patch(ctx, path, opts)
			},
			expected: "PATCH",
		},
		{
			name: "DELETE",
			method: func(ctx context.Context, path string, opts RequestOptions) (*HttpResponse, error) {
				return client.Delete(ctx, path, opts)
			},
			expected: "DELETE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resp, err := tt.method(ctx, "/", RequestOptions{})
			require.NoError(t, err)

			var result testResponse
			err = resp.UnmarshalBody(&result)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.Message)
		})
	}
}

func TestRequestOptions(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		requestOptions RequestOptions
		assertFunc     func(ctx *fasthttp.RequestCtx)
	}{
		{
			name:   "GET with headers",
			method: "GET",
			requestOptions: RequestOptions{
				Header: map[string]string{
					"X-Test-Header": "test-value",
				},
			},
			assertFunc: func(ctx *fasthttp.RequestCtx) {
				assert.Equal(t, "test-value", string(ctx.Request.Header.Peek("X-Test-Header")))
			},
		},
		{
			name:   "GET with query params",
			method: "GET",
			requestOptions: RequestOptions{
				Query: map[string][]string{
					"test_param": {"test-value"},
				},
			},
			assertFunc: func(ctx *fasthttp.RequestCtx) {
				assert.Equal(t, "test-value", string(ctx.QueryArgs().Peek("test_param")))
			},
		},
		{
			name:   "POST with JSON body",
			method: "POST",
			requestOptions: RequestOptions{
				Body: []byte(`{"test_field": "test-value"}`),
			},
			assertFunc: func(ctx *fasthttp.RequestCtx) {
				var body map[string]interface{}
				json.Unmarshal(ctx.PostBody(), &body)
				assert.Equal(t, "test-value", body["test_field"])
			},
		},
		{
			name:   "PUT with form data",
			method: "PUT",
			requestOptions: RequestOptions{
				FormData: map[string][]string{
					"test_form": {"test-value"},
				},
			},
			assertFunc: func(ctx *fasthttp.RequestCtx) {
				assert.Equal(t, "test-value", string(ctx.FormValue("test_form")))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &fasthttp.Server{
				Handler: func(ctx *fasthttp.RequestCtx) {
					tt.assertFunc(ctx)
					ctx.SetContentType("application/json")
					json.NewEncoder(ctx).Encode(testResponse{Message: "success"})
				},
			}

			ln, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			defer ln.Close()

			go server.Serve(ln)

			client, err := New("http://" + ln.Addr().String())
			require.NoError(t, err)

			ctx := context.Background()
			resp, err := client.Do(ctx, tt.method, "/", tt.requestOptions)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode())
		})
	}
}

func TestResponseUnmarshal(t *testing.T) {
	server := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			ctx.SetContentType("application/json")
			json.NewEncoder(ctx).Encode(testResponse{Message: "test message"})
		},
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	go server.Serve(ln)

	client, err := New("http://" + ln.Addr().String())
	require.NoError(t, err)

	ctx := context.Background()
	resp, err := client.Get(ctx, "/", RequestOptions{})
	require.NoError(t, err)

	var result testResponse
	err = resp.UnmarshalBody(&result)
	require.NoError(t, err)
	assert.Equal(t, "test message", result.Message)
}

func TestErrorCases(t *testing.T) {
	// Test with non-existent server
	client, err := New("http://localhost:99999")
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.Get(ctx, "/", RequestOptions{})
	assert.Error(t, err)

	// Test with invalid JSON response
	server := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			ctx.SetContentType("application/json")
			ctx.WriteString("invalid json")
		},
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	go server.Serve(ln)

	client, err = New("http://" + ln.Addr().String())
	require.NoError(t, err)

	resp, err := client.Get(ctx, "/", RequestOptions{})
	require.NoError(t, err)

	var result testResponse
	err = resp.UnmarshalBody(&result)
	assert.Error(t, err)
}

func TestContextTimeout(t *testing.T) {
	server := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			time.Sleep(100 * time.Millisecond)
			ctx.WriteString("{}")
		},
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	go server.Serve(ln)

	client, err := New("http://" + ln.Addr().String())
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err = client.Get(ctx, "/", RequestOptions{})
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}
