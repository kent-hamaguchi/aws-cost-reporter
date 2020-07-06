package notify

import (
	"errors"
	"testing"
	"time"

	"github.com/kent-hamaguchi/aws-cost-reporter/internal/cost"
)

type costRepositoryMock struct {
	getMonthly    cost.Cost
	errGetMonthly error
}

func (c costRepositoryMock) GetMonthly(int, time.Month) (cost.Cost, error) {
	return c.getMonthly, c.errGetMonthly
}

type controllerMock struct {
	sendInput SendInput
}

func (c controllerMock) Send() SendInput {
	return c.sendInput
}

type presenterMock struct {
	errSend error
}

func (p presenterMock) Send(SendOutput) error {
	return p.errSend
}

func Test_notify_Send(t *testing.T) {
	type fields struct {
		costRepo   cost.Repository
		controller Controller
		presenter  Presenter
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				costRepo:   costRepositoryMock{},
				controller: controllerMock{},
				presenter:  presenterMock{},
			},
		},
		{
			name: "repository_error",
			fields: fields{
				costRepo:   costRepositoryMock{errGetMonthly: errors.New("")},
				controller: controllerMock{},
				presenter:  presenterMock{},
			},
			wantErr: true,
		},
		{
			name: "presenter_error",
			fields: fields{
				costRepo:   costRepositoryMock{},
				controller: controllerMock{},
				presenter:  presenterMock{errSend: errors.New("")},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &notify{
				costRepo:   tt.fields.costRepo,
				controller: tt.fields.controller,
				presenter:  tt.fields.presenter,
			}
			if err := r.Send(); (err != nil) != tt.wantErr {
				t.Errorf("notify.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
