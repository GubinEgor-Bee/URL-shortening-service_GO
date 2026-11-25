package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handlePost(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "positive test #1",
			body: "https://practicum.yandex.ru",
			want: want{
				code:        http.StatusCreated,
				response:    `^http://localhost:8080/[a-zA-Z0-9]{8}$`, // адрес начинается с http://localhost:8080
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #2",
			body: "",
			want: want{
				code:        http.StatusBadRequest,
				response:    "invalid url",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			handlePost(w, request)
			res := w.Result()

			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Regexp(t, tt.want.response, string(resBody))

		})
	}
}

func Test_handleGet(t *testing.T) {

	urls = map[string]string{
		"test1": "https://practicum.yandex.ru",
	}

	type want struct {
		code     int
		response string
		location string
	}

	tests := []struct {
		name string // description of this test case
		id   string
		want want
	}{
		{
			name: "positive test #1",
			id:   "test1",
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: "",
				location: "https://practicum.yandex.ru",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			request.SetPathValue("id", tt.id)

			w := httptest.NewRecorder()
			handleGet(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.location, res.Header.Get("Location"))

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			if tt.want.response != "" {
				assert.Regexp(t, tt.want.response, string(resBody))
			}

		})
	}
}
