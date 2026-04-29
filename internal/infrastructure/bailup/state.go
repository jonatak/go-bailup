package bailup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/command"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/model"
)

var ErrDisconnected = errors.New("bailup disconnected")

func (b *Bailup) Execute(ctx context.Context, cmd command.JSONCommand) (*model.State, error) {
	if !b.IsConnected() {
		return nil, NewBailupError("cannot fetch regulation state: client is not connected", ErrDisconnected)
	}

	payload, err := cmd.ToJSON()
	if err != nil {
		return nil, NewBailupError("cannot serialise command", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api-client/regulations/%s", bailupWebsite, b.regulation),
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, NewBailupError("could not build regulation state request", err)
	}

	for k, v := range baseHeader() {
		req.Header.Set(k, v)
	}
	xsrf, err := b.CurrentXSRFToken()
	if err != nil {
		return nil, NewBailupError("could not read current xsrf token", err)
	}

	req.Header.Set("X-Csrf-Token", b.csrf)
	req.Header.Set("X-Xsrf-Token", xsrf)

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, NewBailupError("could not fetch regulation state", errors.Join(err, ErrDisconnected))
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("could not close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, NewBailupError(
			fmt.Sprintf("could not fetch regulation state: unexpected status %d", resp.StatusCode),
			ErrDisconnected,
		)
	}

	decoder := json.NewDecoder(resp.Body)
	var response model.Response
	if err := decoder.Decode(&response); err != nil {
		return nil, NewBailupError("could not decode regulation state response", err)
	}

	return &response.Data, nil
}

func (b *Bailup) GetState(ctx context.Context) (*model.State, error) {
	return b.Execute(ctx, &command.EmptyCommand{})
}
