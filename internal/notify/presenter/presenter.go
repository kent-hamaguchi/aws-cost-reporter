package presenter

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kent-hamaguchi/aws-cost-reporter/internal/notify"
)

// Presenter レポートプレゼンタ
type Presenter struct {
	cfg     Config
	httpCli interface {
		Do(*http.Request) (*http.Response, error)
	}
}

// Config コンフィグ
type Config struct {
	SlackWebhookURL string
}

// New レポートプレゼンタを作成する
func New(cfg Config) *Presenter {
	return &Presenter{
		cfg:     cfg,
		httpCli: &http.Client{},
	}
}

// Send レポートの送信処理
func (p *Presenter) Send(out notify.SendOutput) error {
	s := strconv.FormatFloat(out.Amount, 'f', 2, 64)
	b := `{"text":"` + fmt.Sprintf("AWS 今月分の予想費用: %s $", s) + `"}`
	req, err := http.NewRequest(
		http.MethodPost,
		p.cfg.SlackWebhookURL,
		bytes.NewBuffer([]byte(b)),
	)
	if err != nil {
		return fmt.Errorf("create http request -> %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := p.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("http request -> %w", err)
	}
	res.Body.Close()
	return nil
}
