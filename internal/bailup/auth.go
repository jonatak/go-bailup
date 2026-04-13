package bailup

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func (b *Bailup) Connect() error {
	resp, err := b.client.Get(fmt.Sprintf("%s/client/connexion", bailupWebsite))
	if err != nil {
		return NewBailupError("could not load login page", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return NewBailupError("could not parse login page", err)
	}

	tokenNode := htmlquery.FindOne(doc, "//input[@name=\"_token\"]/@value")
	if tokenNode == nil {
		return NewBailupError("login page did not contain form token", nil)
	}

	csrfNode := htmlquery.FindOne(doc, "//meta[@name=\"csrf-token\"]/@content")
	if csrfNode == nil {
		return NewBailupError("login page did not contain csrf token", nil)
	}

	token := htmlquery.InnerText(tokenNode)
	csrf := htmlquery.InnerText(csrfNode)

	if token == "" {
		return NewBailupError("login form token was empty", nil)
	}
	if csrf == "" {
		return NewBailupError("csrf token was empty", nil)
	}

	b.csrf = csrf
	if err := b.login(token); err != nil {
		return NewBailupError("could not establish authenticated session", err)
	}

	return nil
}

func (b *Bailup) login(token string) error {
	if token == "" {
		return errors.New("login form token was empty")
	}

	form := url.Values{}
	form.Set("email", b.email)
	form.Set("password", b.password)
	form.Set("_token", token)

	resp, err := b.client.Post(
		fmt.Sprintf("%s/client/connexion", bailupWebsite),
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return NewBailupError("could not submit login form", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return NewBailupError(
			fmt.Sprintf("login failed: unexpected status %d", resp.StatusCode),
			nil,
		)
	}

	if resp.Request != nil && resp.Request.URL != nil &&
		resp.Request.URL.Path == "/client/connexion" {
		return errors.New("login failed: still on login page")
	}

	var found bool
	_, found = findCookie(b.client.Jar.Cookies(bailupBaseURL), "XSRF-TOKEN")
	if !found {
		return errors.New("login failed: XSRF token not found")
	}

	return nil
}
