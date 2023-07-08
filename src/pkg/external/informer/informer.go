package informer

import (
	"bytes"
	"encoding/json"
	"github.com/artbred/aliasflux/src/pkg/config"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const (
	SendTelegramMessageEndpoint string = "/telegram/send-message"
)

type (
	SendTelegramMessageRequest struct {
		Token   string `json:"chat_token"`
		Message string `json:"message"`
	}

	CallRequest struct {
		Phone   string `json:"phone"`
		Message string `json:"message"`
	}

	JSONResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	Client struct {
		BaseURL string
	}
)

type Level string

const (
	PaymentsLevel      = Level("payments")
	InternalErrorLevel = Level("internal")
)

var (
	informerTokens = map[Level]string{}
	baseURL        string
)

func SendTelegramMessage(message string, level Level) {
	if config.Debug {
		return
	}

	token, ok := informerTokens[level]
	if !ok {
		logrus.Warningf("Token for level %s is not set", level)
		return
	}

	url := baseURL + SendTelegramMessageEndpoint

	req := SendTelegramMessageRequest{
		Token:   token,
		Message: message,
	}

	b, err := json.Marshal(req)
	if err != nil {
		logrus.Error(err)
		return
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		logrus.Error(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
		return
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
		return
	}

	informerTokens = map[Level]string{
		PaymentsLevel:      os.Getenv("INFORMER_PAYMENTS_TOKEN"),
		InternalErrorLevel: os.Getenv("INFORMER_INTERNAL_TOKEN"),
	}

	baseURL = os.Getenv("INFORMER_BASE_URL")
}
