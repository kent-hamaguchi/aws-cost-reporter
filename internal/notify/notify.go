package notify

import (
	"fmt"
	"time"

	"github.com/kent-hamaguchi/aws-cost-reporter/internal/cost"
)

// Controller コスト通知コントローラ
type Controller interface {
	Send() SendInput
}

// Presenter コスト通知プレゼンタ
type Presenter interface {
	Send(SendOutput) error
}

// Notify コスト通知
type Notify struct {
	costRepo   cost.Repository
	controller Controller
	presenter  Presenter
}

// New コスト通知を作成する
func New(
	costRepo cost.Repository,
	controller Controller,
	presenter Presenter,
) *Notify {
	return &Notify{
		costRepo:   costRepo,
		controller: controller,
		presenter:  presenter,
	}
}

// SendInput コスト通知 送信の入力値
type SendInput struct {
	Now time.Time
}

// SendOutput コスト通知 送信の出力値
type SendOutput struct {
	Amount float64
}

// Send コスト通知を実行する
func (r *Notify) Send() error {
	in := r.controller.Send()
	c, err := r.costRepo.GetMonthly(in.Now.Year(), in.Now.Month())
	if err != nil {
		return fmt.Errorf("get cost -> %w", err)
	}
	out := SendOutput{Amount: c.Amount}
	if err := r.presenter.Send(out); err != nil {
		return fmt.Errorf("send cost -> %w", err)
	}
	return nil
}
