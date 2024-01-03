package endpoint_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shortener/internal/app/endpoint"
	"url-shortener/internal/app/endpoint/mocks"
	"url-shortener/internal/pkg/slogdiscard"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name        string
		alias       string
		targetUrl   string
		aliasLength uint64
		respError   string
		mockError   error
	}{
		{
			name:      "Success",
			alias:     "test_alias",
			targetUrl: "http://www.google.com",
		},
		{
			name:      "Empty alias",
			alias:     "",
			targetUrl: "http://www.google.com",
		},
		{
			name:      "Empty url",
			alias:     "test_alias",
			targetUrl: "",
			respError: "wrong request format",
		},
		{
			name:      "Internal error",
			alias:     "test_alias",
			targetUrl: "http://www.google.com",
			respError: "couldn't save url",
			mockError: errors.New("some internal error"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := fmt.Sprintf(`{"target-url": "%s", "alias": "%s"}`, tc.targetUrl, tc.alias)
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(reqBody)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			repositoryMock := mocks.NewRepository(t)
			if tc.respError == "" || tc.mockError != nil {
				repositoryMock.On("SaveUrl", tc.targetUrl, mock.AnythingOfType("string")).Return(int64(0), tc.mockError).Once()
			}
			unit := endpoint.New(slogdiscard.New(), repositoryMock, tc.aliasLength)
			unit.SaveUrl(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()
			var resp endpoint.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			expectedStatus := "OK"
			if tc.respError != "" {
				expectedStatus = "Error"
			}
			assert.Equal(t, expectedStatus, resp.Status)
			assert.Equal(t, tc.respError, resp.Error)
			if tc.respError == "" {
				if tc.alias != "" {
					assert.Equal(t, tc.alias, resp.Alias)
				} else {
					assert.Equal(t, len(resp.Alias), int(tc.aliasLength))
				}
			}
		})
	}

}

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		targetUrl string
		respError string
		mockError error
	}{
		{
			name:      "Success",
			alias:     "test_alias",
			targetUrl: "http://www.google.com",
		},
		{
			name:      "Empty alias",
			alias:     "",
			targetUrl: "http://www.google.com",
			respError: "empty alias",
		},
		{
			name:      "Url not found for alias",
			alias:     "test_alias",
			targetUrl: "",
			respError: "couldn't find given alias",
			mockError: errors.New("some internal error"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tc.alias), bytes.NewReader([]byte{}))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			repositoryMock := mocks.NewRepository(t)
			if tc.alias != "" {
				repositoryMock.On("GetUrl", tc.alias).Return(tc.targetUrl, tc.mockError).Once()
			}
			unit := endpoint.New(slogdiscard.New(), repositoryMock, 10)
			unit.Redirect(rr, req)

			expectedStatus := http.StatusFound
			if tc.mockError != nil || tc.alias == "" {
				expectedStatus = http.StatusOK
			}
			require.Equal(t, expectedStatus, rr.Code)

			body := rr.Body.String()
			if tc.mockError == nil && tc.alias != "" {
				assert.Equal(t, fmt.Sprintf("<a href=\"%s\">Found</a>.\n\n", tc.targetUrl), body)
			} else {
				var resp endpoint.Response
				assert.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.Equal(t, "Error", resp.Status)
				assert.Equal(t, tc.respError, resp.Error)
			}
		})
	}

}
