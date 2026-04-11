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
	xsrf       string
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
	return b.csrf != "" && b.xsrf != ""
}
