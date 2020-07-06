package presenter

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/kent-hamaguchi/aws-cost-reporter/internal/notify"
)

type httpClientMock struct {
	response *http.Response
	err      error
}

func (h httpClientMock) Do(*http.Request) (*http.Response, error) {
	return h.response, h.err
}

func TestPresenter_Send(t *testing.T) {
	type fields struct {
		cfg     Config
		httpCli interface {
			Do(*http.Request) (*http.Response, error)
		}
	}
	type args struct {
		out notify.SendOutput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				httpCli: httpClientMock{
					response: &http.Response{
						Body: ioutil.NopCloser(bytes.NewReader([]byte{0})),
					},
				},
			},
		},
		{
			name: "http_client_error",
			fields: fields{
				httpCli: httpClientMock{
					err: errors.New(""),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Presenter{
				cfg:     tt.fields.cfg,
				httpCli: tt.fields.httpCli,
			}
			if err := p.Send(tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("Presenter.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
