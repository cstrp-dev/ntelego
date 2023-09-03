package helpers

import (
	"TelegoBot/cmd/config"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"regexp"
	"strings"
)

func New(k, p string) (*IHelper, error) {
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
		key:    k,
		prompt: p,
	}, nil
}

func (h *IHelper) Summarize(text string) (string, error) {
	cfg := config.New()
	safeInput, _ := json.Marshal(cfg.Prompt)
	safeText, _ := json.Marshal(text)
	k, _ := base64.StdEncoding.DecodeString(cfg.OpenAiApiKey)
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

func JSONParse[T any](s string) (T, error) {
	var args T

	if err := json.Unmarshal([]byte(s), &args); err != nil {
		return *(new(T)), err
	}

	return args, nil
}

func CleanUpText(t string) string {
	replacer := regexp.MustCompile("\n{3,}")
	return replacer.ReplaceAllString(t, "\n")
}

func Escape(s string) string {
	var (
		specChars = "-_\\*\\[\\]\\(\\)~`>#\\+=\\|{}\\.!"
		replacer  = regexp.MustCompile("[" + specChars + "]")
	)

	return replacer.ReplaceAllString(s, "\\$0")
}
