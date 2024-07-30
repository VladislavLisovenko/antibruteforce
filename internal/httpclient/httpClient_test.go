package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

var testContent = []byte("test")

func TestHttpClient(t *testing.T) {
	handler := func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/test":
			rw.Write(testContent)
		case "/test3":
			vs := req.URL.Query()
			if vs.Get("testKey") != "test" {
				rw.WriteHeader(http.StatusBadRequest)
			}
		default:
			rw.WriteHeader(http.StatusNotFound)
		}
	}

	ctx := context.Background()
	httpServer := httptest.NewServer(http.HandlerFunc(handler))
	hc := New(httpServer.URL)

	t.Run("Check get", func(t *testing.T) {
		b, err := hc.Get(ctx, "test", nil)
		require.NoError(t, err)
		require.Equal(t, testContent, b)

		_, err = hc.Get(ctx, "test2", nil)
		require.Error(t, err)
	})

	t.Run("Check post", func(t *testing.T) {
		b, err := hc.Post(ctx, "test", nil)
		require.NoError(t, err)
		require.Equal(t, testContent, b)

		_, err = hc.Post(ctx, "test2", nil)
		require.Error(t, err)
	})

	t.Run("Check delete", func(t *testing.T) {
		b, err := hc.Delete(ctx, "test", nil)
		require.NoError(t, err)
		require.Equal(t, testContent, b)

		_, err = hc.Delete(ctx, "test2", nil)
		require.Error(t, err)
	})

	t.Run("Check get params", func(t *testing.T) {
		vs := url.Values{}
		vs.Add("testKey", "test")
		_, err := hc.Get(ctx, "test3", vs)
		require.NoError(t, err)

		vs = url.Values{}
		vs.Add("testKey", "test2")
		_, err = hc.Get(ctx, "test3", vs)
		require.Error(t, err)
	})
}
