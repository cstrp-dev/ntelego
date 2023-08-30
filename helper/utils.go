package helper

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

func (h *IHelper) GetData(key, input string) (string, error) {
	safeInput, _ := json.Marshal(input)
	k, _ := base64.StdEncoding.DecodeString(key)
	var data = strings.NewReader(fmt.Sprintf(`{"model":"gpt-3.5-turbo","messages":[{"role":"user","content":%v}],
	"stream":true}`, string(safeInput)))

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", data)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", string(k))
	resp, err := h.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	type Response struct {
		ID      string `json:"id"`
		Choices []struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
		} `json:"choices"`
	}

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		result := ""
		line := scanner.Text()
		obj := "{}"

		if len(line) > 1 {
			splitLine := strings.Split(line, "data: ")
			if len(splitLine) > 1 {
				obj = splitLine[1]
			}
		}

		var d Response
		if err := json.Unmarshal([]byte(obj), &d); err != nil {
			continue
		}

		if d.Choices != nil {
			result = d.Choices[0].Delta.Content
		}

		fmt.Print(result)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}
