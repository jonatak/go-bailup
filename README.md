# Go-Bailup

Go project to interface with Bailup thermostats.

## Library used

- [htmlquery](https://github.com/antchfx/htmlquery)

## Check connection suggestion

```golang
func (b *Bailup) HasSessionTokens() bool {
	return b.csrf != "" && b.xsrf != ""
}

func (b *Bailup) CheckSession() error {
	req, err := http.NewRequest(http.MethodGet, bailupWebsite+"/client", nil)
	if err != nil {
		return err
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Request != nil && resp.Request.URL != nil &&
		resp.Request.URL.Path == "/client/connexion" {
		return errors.New("not authenticated")
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
```