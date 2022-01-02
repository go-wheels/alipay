# alipay
[![CI](https://github.com/go-wheels/alipay/actions/workflows/ci.yml/badge.svg)](https://github.com/go-wheels/alipay/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-wheels/alipay.svg)](https://pkg.go.dev/github.com/go-wheels/alipay)
[![License](https://img.shields.io/github/license/go-wheels/alipay)](LICENSE)

支付宝 SDK for Go

## 使用指南

安装支付宝 SDK

```shell script
go get -u github.com/go-wheels/alipay
```

初始化支付宝 SDK

```go
alipayOptions := alipay.Options{
    Gateway:         alipay.GatewayProduction, // 支付宝网关
    AppID:           "", // 应用ID
    AppPrivateKey:   "", // 应用私钥
    AlipayPublicKey: "", // 支付宝公钥
}
alipayClient, _ := alipay.NewClient(alipayOptions)
```

## 支付 API

交易查询

```go
tradeQueryRequest := alipay.TradeQueryRequest{
    OutTradeNo: "1586150366616",
}
var tradeQueryResponse alipay.TradeQueryResponse
alipayClient.Execute(tradeQueryRequest, &tradeQueryResponse)
log.Printf("%#v", tradeQueryResponse)
```

条码支付

```go
tradePayRequest := alipay.TradePayRequest{
    NotifyURL:   "http://example.com/notify_url",
    OutTradeNo:  "1586150366616",
    Scene:       "bar_code",
    AuthCode:    "283189274716618278",
    Subject:     "fresh meat",
    TotalAmount: "0.01",
}
var tradePayResponse alipay.TradePayResponse
alipayClient.Execute(tradePayRequest, &tradePayResponse)
log.Printf("%#v", tradePayResponse)
```

扫码支付

```go
tradePrecreateRequest := alipay.TradePrecreateRequest{
    NotifyURL:   "http://example.com/notify_url",
    OutTradeNo:  "1586150366616",
    TotalAmount: "0.01",
    Subject:     "fresh meat",
}
var tradePrecreateResponse alipay.TradePrecreateResponse
alipayClient.Execute(tradePrecreateRequest, &tradePrecreateResponse)
log.Printf("%#v", tradePrecreateResponse)
```

App 支付

```go
tradeAppPayRequest := alipay.TradeAppPayRequest{
    ReturnURL:   "http://example.com/return_url",
    NotifyURL:   "http://example.com/notify_url",
    OutTradeNo:  "1586150366616",
    TotalAmount: "0.01",
    Subject:     "fresh meat",
}
query, _ := alipayClient.SDKExecute(tradeAppPayRequest)
log.Print(query)
```

手机网站支付

```go
tradeWapPayRequest := alipay.TradeWapPayRequest{
    ReturnURL:   "http://example.com/return_url",
    NotifyURL:   "http://example.com/notify_url",
    Subject:     "fresh meat",
    OutTradeNo:  "1586150366616",
    TotalAmount: "0.01",
    QuitURL:     "http://example.com/quit_url",
    ProductCode: "QUICK_WAP_WAY",
}
url, _ := alipayClient.PageExecute(tradeWapPayRequest)
log.Print(url)
```

电脑网站支付

```go
tradePagePayRequest := alipay.TradePagePayRequest{
    ReturnURL:   "http://example.com/return_url",
    NotifyURL:   "http://example.com/notify_url",
    OutTradeNo:  "1586150366616",
    ProductCode: "FAST_INSTANT_TRADE_PAY",
    TotalAmount: "0.01",
    Subject:     "fresh meat",
}
url, _ := alipayClient.PageExecute(tradePagePayRequest)
log.Print(url)
```

通知验签

```go
http.HandleFunc("/notify_url", func(w http.ResponseWriter, r *http.Request) {
    err := alipayClient.VerifyNotification(r)
    if err != nil {
        // 验签失败
        return
    }
    // 验签成功
    w.Write([]byte("success"))
})
```
