package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type GitLabProvider struct {
	*ProviderData
}

func NewGitLabProvider(p *ProviderData) *GitLabProvider {
	p.ProviderName = "GitLab"
	if p.LoginURL == nil || p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{
			Scheme: "https",
			Host:   "gitlab.com",
			Path:   "/oauth/authorize",
		}
	}
	if p.RedeemURL == nil || p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{
			Scheme: "https",
			Host:   "gitlab.com",
			Path:   "/oauth/token",
		}
	}
	if p.ValidateURL == nil || p.ValidateURL.String() == "" {
		p.ValidateURL = &url.URL{
			Scheme: "https",
			Host:   "gitlab.com",
			Path:   "/api/v3/user/emails",
		}
	}
	if p.Scope == "" {
		p.Scope = "api"
	}
	return &GitLabProvider{ProviderData: p}
}

func (p *GitLabProvider) GetEmailAddress(s *SessionState) (string, error) {

	var emails []struct {
		ID 			int   	`json:"id"`
		Email   string	`json:"email"`
	}

	params := url.Values{
		"access_token": {s.AccessToken},
	}
	endpoint := p.ValidateURL.Scheme + "://" + p.ValidateURL.Host + p.ValidateURL.Path + "?" + params.Encode()
	resp, err := http.DefaultClient.Get(endpoint)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("got %d from %q %s", resp.StatusCode, endpoint, body)
	} else {
		log.Printf("got %d from %q %s", resp.StatusCode, endpoint, body)
	}

	if err := json.Unmarshal(body, &emails); err != nil {
		return "", fmt.Errorf("%s unmarshaling %s", err, body)
	}

	for _, email := range emails {
		if email.ID == 1 {
			return email.Email, nil
		}
	}

	return "", nil
}
