package helpers

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"strings"
)

func New() (*IHelper, error) {
	jar := tlsclient.NewCookieJar()
	options := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(120),
		tlsclient.WithClientProfile(tlsclient.Firefox_110),
		tlsclient.WithNotFollowRedirects(),
		tlsclient.WithCookieJar(jar),
	}
	client, err := tlsclient.NewHttpClient(tlsclient.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return &IHelper{
		client: client,
	}, nil
}

func (h *IHelper) GetData(key, prompt, text string) (string, error) {
	safeInput, _ := json.Marshal(prompt)
	safeText, _ := json.Marshal(text)
	k, _ := base64.StdEncoding.DecodeString(key)
	var data = strings.NewReader(fmt.Sprintf(`{"model":"gpt-3.5-turbo","messages":[{"role":"system","content":%v}, {"role":"user","content":%v}],
	"stream":true}`, string(safeInput), string(safeText)))

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", data)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", string(k))
	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	var (
		results []string
		d       Response
	)

	for scanner.Scan() {
		lines := scanner.Text()
		if strings.HasPrefix(lines, "data: ") {
			obj := strings.TrimPrefix(lines, "data: ")
			if err := json.Unmarshal([]byte(obj), &d); err != nil {
				continue
			}

			if d.Choices != nil {
				results = append(results, d.Choices[0].Delta.Content)
			}
		}
	}

	return strings.Join(results, ""), nil
}
