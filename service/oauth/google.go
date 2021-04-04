package oauth

import (
	"github.com/zzokki81/uas/service"
	"golang.org/x/oauth2"
)

type Google struct {
	googleOauthConfig *oauth2.Config
}

func NewGoogle(config *oauth2.Config) *Google {
	return &Google{googleOauthConfig: config}
}

func (g *Google) Oauth() service.GoogleOauth {
	return g.googleOauthConfig
}
