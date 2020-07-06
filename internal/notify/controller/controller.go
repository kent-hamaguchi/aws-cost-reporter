package controller

import (
	"time"

	"github.com/kent-hamaguchi/aws-cost-reporter/internal/notify"
)

// Controller レポートコントローラ
type Controller struct{}

// New レポートコントローラを作成する
func New() *Controller {
	return &Controller{}
}

// Send レポート送信の入力値を返す
func (c *Controller) Send() notify.SendInput {
	return notify.SendInput{
		Now: time.Now(),
	}
}
