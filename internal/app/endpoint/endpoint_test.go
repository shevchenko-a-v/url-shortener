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
			name:      "Empty url",
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
			serviceMock := mocks.NewService(t)
			if tc.respError == "" || tc.mockError != nil {
				serviceMock.On("SaveUrl", tc.targetUrl, mock.AnythingOfType("string")).Return(tc.mockError).Once()
			}
			unit := endpoint.New(slogdiscard.New(), serviceMock, tc.aliasLength)
			unit.SaveUrl(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()
			var resp endpoint.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			expectedStatus := "OK"
			if tc.respError != "" {
				expectedStatus = "Error"
			}
			require.Equal(t, expectedStatus, resp.Status)
			require.Equal(t, tc.respError, resp.Error)
			if tc.respError == "" {
				if tc.alias != "" {
					require.Equal(t, tc.alias, resp.Alias)
				} else {
					require.Equal(t, len(resp.Alias), int(tc.aliasLength))
				}
			}
		})
	}

}
