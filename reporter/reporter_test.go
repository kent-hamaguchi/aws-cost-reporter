package reporter

import (
	"errors"
	"testing"

	"github.com/kent-hamaguchi/aws-cost-reporter/internal/notify"
)

type notifyMock struct {
	err error
}

func (n notifyMock) Send() error {
	return n.err
}

func Test_reporter_Send(t *testing.T) {
	type fields struct {
		notify notify.Notify
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				notify: notifyMock{},
			},
		},
		{
			name: "error",
			fields: fields{
				notify: notifyMock{err: errors.New("")},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reporter{
				notify: tt.fields.notify,
			}
			if err := r.Send(); (err != nil) != tt.wantErr {
				t.Errorf("reporter.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
