package repository

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/costexploreriface"
	"github.com/kent-hamaguchi/aws-cost-reporter/internal/cost"
)

type costExplorerMock struct {
	costexploreriface.ClientAPI
	getCostAndUsageRequest costexplorer.GetCostAndUsageRequest
}

func (c costExplorerMock) GetCostAndUsageRequest(*costexplorer.GetCostAndUsageInput) costexplorer.GetCostAndUsageRequest {
	return c.getCostAndUsageRequest
}

func TestRepository_GetMonthly(t *testing.T) {
	type fields struct {
		ce costexploreriface.ClientAPI
	}
	type args struct {
		y int
		m time.Month
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    cost.Cost
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				ce: costExplorerMock{
					getCostAndUsageRequest: costexplorer.GetCostAndUsageRequest{
						Request: &aws.Request{
							Data: &costexplorer.GetCostAndUsageOutput{
								ResultsByTime: []costexplorer.ResultByTime{
									{
										Groups: []costexplorer.Group{
											{
												Metrics: map[string]costexplorer.MetricValue{
													metrics: {
														Amount: aws.String("1.5"),
													},
												},
											},
											{
												Metrics: map[string]costexplorer.MetricValue{
													metrics: {
														Amount: aws.String("2.0"),
													},
												},
											},
										},
									},
								},
							},
							Error:       nil,
							HTTPRequest: &http.Request{},
							Retryer:     aws.NoOpRetryer{},
						},
					},
				},
			},
			want: cost.Cost{
				Amount: 3.5,
			},
		},
		{
			name: "costexplorer_error",
			fields: fields{
				ce: costExplorerMock{
					getCostAndUsageRequest: costexplorer.GetCostAndUsageRequest{
						Request: &aws.Request{
							Error:       errors.New(""),
							HTTPRequest: &http.Request{},
							Retryer:     aws.NoOpRetryer{},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "amount_is_not_float",
			fields: fields{
				ce: costExplorerMock{
					getCostAndUsageRequest: costexplorer.GetCostAndUsageRequest{
						Request: &aws.Request{
							Data: &costexplorer.GetCostAndUsageOutput{
								ResultsByTime: []costexplorer.ResultByTime{
									{
										Groups: []costexplorer.Group{
											{
												Metrics: map[string]costexplorer.MetricValue{
													metrics: {
														Amount: aws.String("a"),
													},
												},
											},
										},
									},
								},
							},
							Error:       nil,
							HTTPRequest: &http.Request{},
							Retryer:     aws.NoOpRetryer{},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				ce: tt.fields.ce,
			}
			got, err := r.GetMonthly(tt.args.y, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetMonthly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetMonthly() = %v, want %v", got, tt.want)
			}
		})
	}
}
