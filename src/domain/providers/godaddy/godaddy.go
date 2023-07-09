package godaddy

import (
	"encoding/json"
	"fmt"
	"github.com/artbred/aliasflux/src/pkg/config"
	"io"
	"net/http"
)

type GoDaddy struct {
	BaseURL   string
	APIKey    string
	APISecret string
}

type Tld struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (g *GoDaddy) ListTld() (tlds []Tld, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/domains/tlds", g.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", g.APIKey, g.APISecret))

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err = json.Unmarshal(body, &tlds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return
}

func NewClient() *GoDaddy {
	domain := "https://api.godaddy.com"
	//if config.Debug {
	//	domain = "https://api.ote-godaddy.com"
	//}

	return &GoDaddy{
		APIKey:    config.GoDaddyApiKey,
		APISecret: config.GoDaddyApiSecret,
		BaseURL:   domain,
	}
}
