package bailup

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jonatak/go-bailup/internal/bailup/model"
)

func (b *Bailup) GetState() (*model.State, error) {
	if !b.IsConnected() {
		return nil, NewBailupError("cannot fetch regulation state: client is not connected", nil)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api-client/regulations/%s", bailupWebsite, b.regulation),
		nil,
	)
	if err != nil {
		return nil, NewBailupError("could not build regulation state request", err)
	}

	for k, v := range baseHeader() {
		req.Header.Set(k, v)
	}

	req.Header.Set("X-Csrf-Token", b.csrf)
	req.Header.Set("X-Xsrf-Token", b.xsrf)

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, NewBailupError("could not fetch regulation state", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewBailupError(
			fmt.Sprintf("could not fetch regulation state: unexpected status %d", resp.StatusCode),
			nil,
		)
	}

	decoder := json.NewDecoder(resp.Body)
	var response model.Response
	if err := decoder.Decode(&response); err != nil {
		return nil, NewBailupError("could not decode regulation state response", err)
	}

	return &response.Data, nil
}
