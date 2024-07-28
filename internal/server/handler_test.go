package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/VladislavLisovenko/antibruteforce/internal/app"
	"github.com/VladislavLisovenko/antibruteforce/internal/list"
	"github.com/VladislavLisovenko/antibruteforce/internal/logger"
	"github.com/VladislavLisovenko/antibruteforce/internal/ratelimit"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

const (
	wlKey           = "whiteList"
	blKey           = "blackList"
	loggerLevel     = "debug"
	responseTrue    = "{\"ok\":true}"
	responseFalse   = "{\"ok\":false}"
	testWhiteListIP = "127.0.0.3/32"
	testBlackListIP = "127.0.0.4/32"
)

type structTestData struct {
	login      string
	password   string
	ip         string
	statusCode int
	response   string
}

var (
	checkTestData = []structTestData{
		{
			"test",
			"pass",
			"127.0.0.1",
			http.StatusOK,
			responseTrue,
		},
		{
			"test",
			"pass",
			"127.0.0.1",
			http.StatusOK,
			responseTrue,
		},
		{
			"test",
			"pass",
			"127.0.0.1",
			http.StatusOK,
			responseFalse,
		},
	}
	resetTestData = []structTestData{
		{
			"test2",
			"pass2",
			"127.0.0.2",
			http.StatusOK,
			responseTrue,
		},
		{
			"test2",
			"pass2",
			"127.0.0.2",
			http.StatusOK,
			responseTrue,
		},
		{
			"test2",
			"pass2",
			"127.0.0.2",
			http.StatusOK,
			responseFalse,
		},
	}
	whiteListTestData = []structTestData{
		{
			"test3",
			"pass3",
			"127.0.0.3",
			http.StatusOK,
			responseTrue,
		},
		{
			"test3",
			"pass3",
			"127.0.0.3",
			http.StatusOK,
			responseTrue,
		},
		{
			"test3",
			"pass3",
			"127.0.0.3",
			http.StatusOK,
			responseTrue,
		},
	}
	blackListTestData = []structTestData{
		{
			"test4",
			"pass4",
			"127.0.0.4",
			http.StatusOK,
			responseFalse,
		},
	}
)

func TestHandler(t *testing.T) {
	s := miniredis.RunT(t)
	r := redis.NewClient(&redis.Options{Addr: s.Addr()})
	ctx := context.Background()
	rt := ratelimit.New(2, 50, 20, 50, 10)
	wl, err := list.New(ctx, wlKey, *r)
	require.NoError(t, err)
	bl, err := list.New(ctx, blKey, *r)
	require.NoError(t, err)
	a := app.NewApp(rt, wl, bl)
	l, err := logger.New(loggerLevel)
	require.NoError(t, err)

	h := NewHandlers(a, l)

	httpClient := &http.Client{}

	t.Run("Page not found", func(t *testing.T) {
		test := httptest.NewServer(h.Handlers(ctx))
		defer test.Close()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, test.URL, nil)
		require.NoError(t, err)
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Check", func(t *testing.T) {
		testServer := httptest.NewServer(h.Handlers(ctx))
		defer testServer.Close()
		for _, testData := range checkTestData {
			check(ctx, t, testServer, testData)
		}
	})

	t.Run("Reset", func(t *testing.T) {
		testServer := httptest.NewServer(h.Handlers(ctx))
		defer testServer.Close()
		for _, testData := range resetTestData {
			check(ctx, t, testServer, testData)
		}

		vs := url.Values{}
		vs.Add(LoginField, "test2")
		vs.Add(LoginField, "127.0.0.2")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServer.URL+"/reset?"+vs.Encode(), nil)
		require.NoError(t, err)
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		for _, testData := range resetTestData {
			check(ctx, t, testServer, testData)
		}
	})

	t.Run("Whitelist", func(t *testing.T) {
		testServer := httptest.NewServer(h.Handlers(ctx))
		defer testServer.Close()

		vs := url.Values{}
		vs.Add(IPField, testWhiteListIP)
		addToList(ctx, t, testServer, "whiteList", vs)

		for _, testData := range whiteListTestData {
			check(ctx, t, testServer, testData)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, testServer.URL+"/whiteList?"+vs.Encode(), nil)
		require.NoError(t, err)
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("BlackList", func(t *testing.T) {
		testServer := httptest.NewServer(h.Handlers(ctx))
		defer testServer.Close()

		vs := url.Values{}
		vs.Add(IPField, testBlackListIP)
		addToList(ctx, t, testServer, "blackList", vs)

		for _, testData := range blackListTestData {
			check(ctx, t, testServer, testData)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, testServer.URL+"/blackList?"+vs.Encode(), nil)
		require.NoError(t, err)
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Not allowed", func(t *testing.T) {
		test := httptest.NewServer(h.Handlers(ctx))
		defer test.Close()

		req, err := http.NewRequestWithContext(ctx, http.MethodPut, test.URL+"/blackList", nil)
		require.NoError(t, err)
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})
}

func check(ctx context.Context, t *testing.T, test *httptest.Server, testData structTestData) {
	t.Helper()
	httpClient := &http.Client{}

	vs := url.Values{}
	vs.Add(LoginField, testData.login)
	vs.Add(PasswordField, testData.password)
	vs.Add(IPField, testData.ip)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, test.URL+"/check?"+vs.Encode(), nil)
	require.NoError(t, err)
	resp, err := httpClient.Do(req)

	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, testData.statusCode, resp.StatusCode)

	out, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, testData.response, string(out))
}

func addToList(ctx context.Context, t *testing.T, test *httptest.Server, typeList string, data url.Values) {
	t.Helper()
	httpClient := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, test.URL+"/"+typeList+"?"+data.Encode(), nil)
	require.NoError(t, err)
	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
