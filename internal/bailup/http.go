package bailup

import "net/http"

const bailupWebsite = "https://www.baillconnect.com"

func findCookie(cookies []*http.Cookie, name string) (string, bool) {
	for _, c := range cookies {
		if c.Name == name {
			return c.Value, true
		}
	}
	return "", false
}

func baseHeader() map[string]string {
	return map[string]string{
		"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
		"Accept":             "application/json, text/plain, */*",
		"Accept-Encoding":    "gzip, deflate, br, zstd",
		"Accept-Language":    "fr-FR,fr;q=0.9",
		"Content-Type":       "application/json;charset=UTF-8",
		"X-Requested-With":   "XMLHttpRequest",
		"Sec-Ch-Ua":          "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Google Chrome\";v=\"128\"",
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": "\"macOS\"",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
	}
}
