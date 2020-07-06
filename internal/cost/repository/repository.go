package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/costexploreriface"
	"github.com/kent-hamaguchi/aws-cost-reporter/internal/cost"
)

const (
	dateFmt      = "2006-01-02"
	metrics      = "UnblendedCost"
	dimentionKey = "LINKED_ACCOUNT"
)

// Repository コストリポジトリ
type Repository struct {
	ce costexploreriface.ClientAPI
}

// New コストリポジトリを作成する
func New(c aws.Config) *Repository {
	return &Repository{
		ce: costexplorer.New(c),
	}
}

// GetMonthly 月次コストを取得する
func (r *Repository) GetMonthly(y int, m time.Month) (cost.Cost, error) {
	s := time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
	e := s.AddDate(0, 1, 0)
	req := r.ce.GetCostAndUsageRequest(&costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(s.Format(dateFmt)),
			End:   aws.String(e.Format(dateFmt)),
		},
		Granularity: costexplorer.GranularityMonthly,
		Metrics:     []string{metrics},
		GroupBy: []costexplorer.GroupDefinition{
			{
				Type: costexplorer.GroupDefinitionTypeDimension,
				Key:  aws.String(dimentionKey),
			},
		},
	})
	res, err := req.Send(context.Background())
	if err != nil {
		return cost.Cost{}, fmt.Errorf("AWS CostExplorer get cost and usage request -> %w", err)
	}
	out := res.GetCostAndUsageOutput
	var amount float64
	for _, r := range out.ResultsByTime {
		for _, g := range r.Groups {
			f, err := strconv.ParseFloat(*g.Metrics[metrics].Amount, 64)
			if err != nil {
				return cost.Cost{}, fmt.Errorf("AWS CostExplorer returns illegal amount -> %w", err)
			}
			amount += f
		}
	}
	return cost.Cost{
		Amount: amount,
	}, nil
}
