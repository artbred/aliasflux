package godaddy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/artbred/aliasflux/src/pkg/config"
	"io"
	"net/http"
	"net/url"
	"time"
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

type NameDomains struct {
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
}

func (g *GoDaddy) SuggestDomains(query string, tlds string) (nameDomains *NameDomains, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	baseURL, err := url.Parse(fmt.Sprintf("%s/v1/domains/suggest", g.BaseURL))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("query", query)
	params.Add("tlds", tlds)
	params.Add("limit", "5")
	params.Add("waitMs", "1000")

	// Add parameters to the URL
	baseURL.RawQuery = params.Encode()

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", g.APIKey, g.APISecret))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	nameDomains = &NameDomains{
		Name: query,
	}

	var domains = []map[string]string{}

	err = json.Unmarshal(body, &domains)

	for _, domain := range domains {
		for _, value := range domain {
			nameDomains.Domains = append(nameDomains.Domains, value)
		}
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
