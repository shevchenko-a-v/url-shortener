package tests

import (
	"net/http"
	"net/url"
	"testing"
	"url-shortener/internal/app/endpoint"
	"url-shortener/internal/pkg/alias"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

const host = "localhost:8082"

func TestUrlShortener_SmokeTest(t *testing.T) {
	u := url.URL{Scheme: "http", Host: host}
	e := httpexpect.Default(t, u.String())

	e.POST("/save").WithJSON(endpoint.Request{TargetUrl: gofakeit.URL(), Alias: alias.CreateRandom(10)}).
		WithBasicAuth("myuser", "mypass").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("alias")
}
