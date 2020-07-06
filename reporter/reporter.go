package reporter

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kent-hamaguchi/aws-cost-reporter/internal/cost/repository"
	"github.com/kent-hamaguchi/aws-cost-reporter/internal/notify"
	"github.com/kent-hamaguchi/aws-cost-reporter/internal/notify/controller"
	"github.com/kent-hamaguchi/aws-cost-reporter/internal/notify/presenter"
)

// Reporter AWSのコストをSlackに通知します
type Reporter interface {
	Send() error
}

// Config コンフィグ
type Config struct {
	// AWSConfig AWS認証情報
	AWSConfig aws.Config
	// SlackWebhookURL 通知先のSlackIncomingWebhookのURL
	SlackWebhookURL string
}

type reporter struct {
	notify notify.Notify
}

// New コスト通知のインスタンスを作成します
func New(cfg Config) Reporter {
	r := repository.New(cfg.AWSConfig)
	c := controller.New()
	p := presenter.New(presenter.Config{SlackWebhookURL: cfg.SlackWebhookURL})
	return &reporter{
		notify: notify.New(r, c, p),
	}
}

// Send コスト通知を送信します
func (r *reporter) Send() error {
	if err := r.notify.Send(); err != nil {
		return fmt.Errorf("send report error -> %w", err)
	}
	return nil
}
