package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zzokki81/uas/service"
)

var (
	oauthStateString = "random-string"
	oauthGetUserInfo = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type Google interface {
	Oauth() service.GoogleOauth
}

type GoogleHandler struct {
	google Google
}

func NewGoogleHandler(google Google) *GoogleHandler {
	return &GoogleHandler{google: google}
}

func (gh *GoogleHandler) HandleMain(c echo.Context) error {
	url := gh.google.Oauth().AuthCodeURL(oauthStateString)
	html := `<html>
<body>
	<a href="` + url + `">Google Log In</a>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}

func (gh *GoogleHandler) ExchangeAuthCodeForTOken(c echo.Context) error {
	if c.QueryParam("state") != oauthStateString {
		return errors.New("State did not match")
	}

	oauthToken, err := gh.google.Oauth().Exchange(c.Request().Context(), c.QueryParam("code"))
	if err != nil {
		return err
	}

	response, err := http.Get(oauthGetUserInfo + oauthToken.AccessToken)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.Response().Writer, "UserInfo: %s\n", contents)

	return nil
}
