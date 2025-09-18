package helpers

import (
	"TelegoBot/cmd/config"
	"bufio"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/sirupsen/logrus"
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

	logrus.Infof("AI Summarization - Input text length: %d chars, Prompt: %s", len(text), cfg.Prompt[:100]+"...")

	var data = strings.NewReader(fmt.Sprintf(`{"model":"gpt-3.5-turbo","messages":[{"role":"system","content":%v}, {"role":"user","content":%v}],
	"stream":true}`, string(safeInput), string(safeText)))

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", data)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", string("Bearer "+cfg.OpenAiApiKey))
	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logrus.Errorf("OpenAI API error - Status: %d", resp.StatusCode)
		return "", fmt.Errorf("OpenAI API returned status %d", resp.StatusCode)
	}

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

	finalSummary := strings.Join(results, "")
	logrus.Infof("AI Summarization completed - Output length: %d chars", len(finalSummary))

	if len(finalSummary) < 50 {
		logrus.Warnf("AI summary too short, possible issue with API response")
	}

	return finalSummary, nil
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
