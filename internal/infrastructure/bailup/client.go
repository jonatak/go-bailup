package bailup

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

type Bailup struct {
	email      string
	password   string
	regulation string
	csrf       string
	client     *http.Client
}

func NewBailup(email, password, regulation string) *Bailup {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create cookie jar: %v", err))
	}

	return &Bailup{
		email:      email,
		password:   password,
		regulation: regulation,
		client:     &http.Client{Jar: jar},
	}
}

func (b *Bailup) IsConnected() bool {
	xsrf, err := b.CurrentXSRFToken()
	if err != nil {
		return false
	}

	return b.csrf != "" && xsrf != ""
}

func (b *Bailup) CurrentXSRFToken() (string, error) {
	xsrf, found := findCookie(b.client.Jar.Cookies(bailupBaseURL), "XSRF-TOKEN")
	if !found {
		return "", NewBailupError("xsrf token not found", ErrDisconnected)
	}

	return xsrf, nil
}
