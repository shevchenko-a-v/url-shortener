package endpoint

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"url-shortener/internal/app/middleware"
	"url-shortener/internal/pkg/alias"
	"url-shortener/internal/pkg/api/response"

	"github.com/go-playground/validator/v10"
)

//go:generate go run github.com/vektra/mockery/v2@v2.38.0 --name=Repository
type Repository interface {
	SaveUrl(string, string) (int64, error)
	GetUrl(string) (string, error)
}

type Request struct {
	TargetUrl string `json:"target-url" validate:"required,url"`
	Alias     string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type Endpoint struct {
	repository         Repository
	log                *slog.Logger
	defaultAliasLength uint64
}

func New(log *slog.Logger, repository Repository, defaultAliasLength uint64) *Endpoint {
	return &Endpoint{repository: repository, log: log, defaultAliasLength: defaultAliasLength}
}

func (e *Endpoint) SaveUrl(rw http.ResponseWriter, req *http.Request) {
	encoder := json.NewEncoder(rw)
	log := e.log.With(
		slog.String("op", "Endpoint.SaveUrl"),
		slog.String("request_id", middleware.GetRequestId(req.Context())),
	)
	var body Request
	if req.Method != http.MethodPost {
		msg := fmt.Sprintf("Not supported request type: %s", req.Method)
		log.Info(msg)
		_ = encoder.Encode(response.Error(msg))
		return
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(fmt.Sprintf("Wrong request: %s", err.Error()))
		_ = encoder.Encode(response.Error("wrong request format"))
		return
	}
	if err = validator.New().Struct(body); err != nil {
		log.Error(fmt.Sprintf("invalid request: %s", err.Error()))
		_ = encoder.Encode(response.Error("wrong request format"))
		return
	}
	if body.TargetUrl == "" {
		log.Error("field url must not be empty")
		_ = encoder.Encode(response.Error("wrong request format"))
		return
	}
	if body.Alias == "" {
		body.Alias = alias.CreateRandom(e.defaultAliasLength)
	}
	var id int64
	id, err = e.repository.SaveUrl(body.TargetUrl, body.Alias)
	if err != nil {
		log.Error(err.Error())
		_ = encoder.Encode(response.Error("couldn't save url"))
		return
	}
	log.Debug(fmt.Sprintf("saved with id: %d", id))
	_ = encoder.Encode(Response{Response: response.OK(), Alias: body.Alias})
}

func (e *Endpoint) Redirect(rw http.ResponseWriter, req *http.Request) {
	encoder := json.NewEncoder(rw)
	log := e.log.With(
		slog.String("op", "Endpoint.Redirect"),
		slog.String("request_id", middleware.GetRequestId(req.Context())),
	)
	if req.Method != http.MethodGet {
		msg := fmt.Sprintf("Not supported request type: %s", req.Method)
		log.Info(msg)
		_ = encoder.Encode(response.Error(msg))
		return
	}
	aliasToGetUrl, _ := strings.CutPrefix(req.URL.String(), "/")
	if aliasToGetUrl == "" {
		log.Error("alias is empty")
		_ = encoder.Encode(response.Error("empty alias"))
		return
	}
	resultUrl, err := e.repository.GetUrl(aliasToGetUrl)
	if err != nil {
		log.Error(err.Error())
		_ = encoder.Encode(response.Error("couldn't find given alias"))
		return
	}
	log.Debug(fmt.Sprintf("found %s url for given alias %s", resultUrl, aliasToGetUrl))

	http.Redirect(rw, req, resultUrl, http.StatusFound)
}
