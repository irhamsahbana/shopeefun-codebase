package integration

import (
	"codebase-app/internal/integration/oauth2google/entity"
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Oauth2googleContract interface {
	GetUrl(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (entity.UserInfoResponse, error)
}

type ouath2google struct {
	cfg oauth2.Config
}

func NewOauth2googleIntegration(clientID, clientSecret, redirectURL string) Oauth2googleContract {
	var googleOauthCfg = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return &ouath2google{
		cfg: googleOauthCfg,
	}
}

func (o *ouath2google) GetUrl(state string, opts ...oauth2.AuthCodeOption) string {
	return o.cfg.AuthCodeURL("state", opts...)
}

func (o *ouath2google) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.cfg.Exchange(ctx, code)
}

func (o *ouath2google) GetUserInfo(ctx context.Context, token *oauth2.Token) (entity.UserInfoResponse, error) {
	var (
		info   = entity.UserInfoResponse{}
		client = o.cfg.Client(ctx, token)
	)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return info, err
	}

	return info, nil
}
