// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package admin_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/testcontext"
)

func assertGet(ctx context.Context, t *testing.T, link string, expected string, authToken string) {
	t.Helper()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", authToken)

	response, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	data, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())

	require.Equal(t, http.StatusOK, response.StatusCode, string(data))
	require.Equal(t, expected, string(data))
}

// assertReq asserts the request and it's OK it returns the response body.
func assertReq(
	ctx *testcontext.Context, t *testing.T, link string, method string, body string,
	expectedStatus int, expectedBody string, authToken string,
) []byte {
	t.Helper()

	var (
		req *http.Request
		err error
	)
	if body == "" {
		req, err = http.NewRequestWithContext(ctx, method, link, nil)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, link, strings.NewReader(body))
	}
	require.NoError(t, err)

	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req) //nolint:bodyclose
	require.NoError(t, err)
	defer ctx.Check(res.Body.Close)

	require.Equal(t, expectedStatus, res.StatusCode, "response status code")

	resBody, err := ioutil.ReadAll(res.Body)
	if expectedBody != "" {
		require.NoError(t, err)
		require.Equal(t, expectedBody, string(resBody), "response body")
	}

	return resBody
}
