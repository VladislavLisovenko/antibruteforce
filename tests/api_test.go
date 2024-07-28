package tests

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type apiTest struct {
	responseStatusCode int
	responseBody       []byte
}

func (t *apiTest) iSendRequestTo(method string, addr string) error {
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, method, addr, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(addr)
	defer resp.Body.Close()
	t.responseStatusCode = resp.StatusCode
	t.responseBody, err = io.ReadAll(resp.Body)

	return err
}

func (t *apiTest) theResponseCodeShouldBe(code int) error {
	if t.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", t.responseStatusCode, code)
	}
	return nil
}

func (t *apiTest) theResponseShouldMatchText(text string) error {
	if string(t.responseBody) != text {
		return fmt.Errorf("unexpected text: %s != %s", t.responseBody, text)
	}
	return nil
}
