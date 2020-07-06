# aws-cost-reporter

![Test](https://github.com/kent-hamaguchi/aws-cost-reporter/workflows/Test/badge.svg)

AWSのCostExplorerから対象アカウントの当月分のコストを取得し、Slackに通知するためのGoのプログラムです。

とりあえずサマリの数値を通知するだけですが、最終的にはAWS Organizationsに属するアカウントすべてのコストを通知できるようにするため、内部的にはそれに繋げていくためのコードが書かれています。

## 始め方

go getで追加。

```bash
go get github.com/kent-hamaguchi/aws-cost-reporter
```

起動コード。Configに渡す値をそれぞれ作成して実行してください。

```go
r := reporter.New(reporter.Config{
    AWSConfig:       cfg,
    SlackWebhookURL: webhook,
})
err := r.Send()
```
